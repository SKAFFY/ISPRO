package logs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	defaultEndpoint = "http://localhost:9428"
	insertPath      = "/insert/jsonline?_stream_fields=level,service&_time_field=time&_msg_field=msg"
)

type VictoriaLogsHandler struct {
	client   *http.Client
	endpoint string
	service  string
	level    slog.Leveler
}

func NewVictoriaLogsHandler(endpoint, service string, level slog.Leveler) *VictoriaLogsHandler {
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	return &VictoriaLogsHandler{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		endpoint: endpoint,
		service:  service,
		level:    level,
	}
}

func (h *VictoriaLogsHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *VictoriaLogsHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := make(map[string]any, record.NumAttrs()+4)

	attrs["service"] = h.service
	attrs["level"] = record.Level.String()
	attrs["time"] = record.Time.Format(time.RFC3339Nano)
	attrs["msg"] = record.Message

	record.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	body, err := json.Marshal(attrs)
	if err != nil {
		return fmt.Errorf("marshal log: %w", err)
	}

	_, _ = os.Stdout.Write(append(body, '\n'))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint+insertPath, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("send log to victorialogs: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (h *VictoriaLogsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *VictoriaLogsHandler) WithGroup(name string) slog.Handler {
	return h
}
