package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultEndpoint = "http://localhost:3100"
	pushPath        = "/loki/api/v1/push"
)

type lokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type lokiPushRequest struct {
	Streams []lokiStream `json:"streams"`
}

type LokiHandler struct {
	client   *http.Client
	endpoint string
	service  string
	level    slog.Leveler
	buffer   map[string]*lokiStream
}

func NewLokiHandler(endpoint, service string, level slog.Leveler) *LokiHandler {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	return &LokiHandler{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		endpoint: endpoint,
		service:  service,
		level:    level,
		buffer:   make(map[string]*lokiStream),
	}
}

func (h *LokiHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *LokiHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := make(map[string]any, record.NumAttrs()+2)
	record.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	logLine := map[string]any{
		"service": h.service,
		"level":   record.Level.String(),
		"msg":     record.Message,
	}
	for k, v := range attrs {
		logLine[k] = v
	}

	line, err := json.Marshal(logLine)
	if err != nil {
		return fmt.Errorf("marshal log: %w", err)
	}

	_, _ = os.Stdout.Write(append(line, '\n'))

	if err := h.sendToLoki(ctx, record, line); err != nil {
		return fmt.Errorf("send log to loki: %w", err)
	}

	return nil
}

func (h *LokiHandler) sendToLoki(ctx context.Context, record slog.Record, line []byte) error {
	stream := map[string]string{
		"service": h.service,
		"level":   record.Level.String(),
	}

	ts := strconv.FormatInt(record.Time.UnixNano(), 10)

	req := lokiPushRequest{
		Streams: []lokiStream{
			{
				Stream: stream,
				Values: [][]string{{ts, string(line)}},
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal push request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint+pushPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("loki push rejected: %s", resp.Status)
	}

	return nil
}

func (h *LokiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *LokiHandler) WithGroup(name string) slog.Handler {
	return h
}
