package entry

import (
	"log/slog"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/mkheyfets/ispro-app/internal/entities/entry/domain"
	"github.com/mkheyfets/ispro-app/internal/models"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations/entries"
)

type Handler struct {
	listEntries *ListEntriesUseCase
	getEntry    *GetEntryUseCase
	createEntry *CreateEntryUseCase
	updateEntry *UpdateEntryUseCase
	deleteEntry *DeleteEntryUseCase
}

func NewHandler(
	listEntries *ListEntriesUseCase,
	getEntry *GetEntryUseCase,
	createEntry *CreateEntryUseCase,
	updateEntry *UpdateEntryUseCase,
	deleteEntry *DeleteEntryUseCase,
) *Handler {
	return &Handler{
		listEntries: listEntries,
		getEntry:    getEntry,
		createEntry: createEntry,
		updateEntry: updateEntry,
		deleteEntry: deleteEntry,
	}
}

func (h *Handler) ListEntries(params entries.ListEntriesParams) middleware.Responder {
	result, err := h.listEntries.Handle(params.HTTPRequest.Context())
	if err != nil {
		slog.Error("error listing entries", "error", err)
		return entries.NewListEntriesOK()
	}
	slog.Info("entries listed", "count", len(result))
	return entries.NewListEntriesOK().WithPayload(toModelList(result))
}

func (h *Handler) CreateEntry(params entries.CreateEntryParams) middleware.Responder {
	cmd := CreateEntryCommand{
		Title:   *params.Body.Title,
		Content: params.Body.Content,
	}
	result, err := h.createEntry.Handle(params.HTTPRequest.Context(), cmd)
	if err != nil {
		slog.Warn("error creating entry", "error", err, "title", cmd.Title)
		return entries.NewCreateEntryCreated()
	}
	slog.Info("entry created", "id", result.ID, "title", result.Title)
	return entries.NewCreateEntryCreated().WithPayload(toModel(result))
}

func (h *Handler) GetEntry(params entries.GetEntryParams) middleware.Responder {
	query := GetEntryQuery{ID: params.ID}
	result, err := h.getEntry.Handle(params.HTTPRequest.Context(), query)
	if err != nil {
		slog.Warn("entry not found", "id", params.ID)
		return entries.NewGetEntryNotFound()
	}
	slog.Info("entry retrieved", "id", result.ID)
	return entries.NewGetEntryOK().WithPayload(toModel(result))
}

func (h *Handler) UpdateEntry(params entries.UpdateEntryParams) middleware.Responder {
	cmd := UpdateEntryCommand{
		ID:      params.ID,
		Title:   *params.Body.Title,
		Content: params.Body.Content,
	}
	result, err := h.updateEntry.Handle(params.HTTPRequest.Context(), cmd)
	if err != nil {
		slog.Warn("error updating entry", "error", err, "id", cmd.ID)
		return entries.NewUpdateEntryNotFound()
	}
	slog.Info("entry updated", "id", result.ID)
	return entries.NewUpdateEntryOK().WithPayload(toModel(result))
}

func (h *Handler) DeleteEntry(params entries.DeleteEntryParams) middleware.Responder {
	cmd := DeleteEntryCommand{ID: params.ID}
	err := h.deleteEntry.Handle(params.HTTPRequest.Context(), cmd)
	if err != nil {
		slog.Warn("error deleting entry", "error", err, "id", cmd.ID)
		return entries.NewDeleteEntryNotFound()
	}
	slog.Info("entry deleted", "id", cmd.ID)
	return entries.NewDeleteEntryNoContent()
}

func toModel(e *domain.Entry) *models.Entry {
	if e == nil {
		return nil
	}
	return &models.Entry{
		ID:        e.ID,
		Title:     e.Title,
		Content:   e.Content,
		CreatedAt: strfmt.DateTime(e.CreatedAt),
		UpdatedAt: strfmt.DateTime(e.UpdatedAt),
	}
}

func toModelList(entries []*domain.Entry) []*models.Entry {
	result := make([]*models.Entry, len(entries))
	for i, e := range entries {
		result[i] = toModel(e)
	}
	return result
}
