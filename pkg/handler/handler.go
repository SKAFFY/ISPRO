package handler

import (
	"log"

	"github.com/go-openapi/runtime/middleware"

	"github.com/mkheyfets/ispro-app/pkg/service"
	"github.com/mkheyfets/ispro-app/restapi/operations/entries"
	"github.com/mkheyfets/ispro-app/restapi/operations/links"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListEntries(params entries.ListEntriesParams) middleware.Responder {
	result, err := h.svc.ListEntries(params.HTTPRequest.Context())
	if err != nil {
		log.Printf("error listing entries: %v", err)
		return entries.NewListEntriesOK()
	}
	return entries.NewListEntriesOK().WithPayload(result)
}

func (h *Handler) CreateEntry(params entries.CreateEntryParams) middleware.Responder {
	result, err := h.svc.CreateEntry(params.HTTPRequest.Context(), *params.Body.Title, params.Body.Content)
	if err != nil {
		log.Printf("error creating entry: %v", err)
		return entries.NewCreateEntryCreated()
	}
	return entries.NewCreateEntryCreated().WithPayload(result)
}

func (h *Handler) GetEntry(params entries.GetEntryParams) middleware.Responder {
	result, err := h.svc.GetEntry(params.HTTPRequest.Context(), params.ID)
	if err != nil {
		return entries.NewGetEntryNotFound()
	}
	return entries.NewGetEntryOK().WithPayload(result)
}

func (h *Handler) UpdateEntry(params entries.UpdateEntryParams) middleware.Responder {
	result, err := h.svc.UpdateEntry(params.HTTPRequest.Context(), params.ID, *params.Body.Title, params.Body.Content)
	if err != nil {
		return entries.NewUpdateEntryNotFound()
	}
	return entries.NewUpdateEntryOK().WithPayload(result)
}

func (h *Handler) DeleteEntry(params entries.DeleteEntryParams) middleware.Responder {
	err := h.svc.DeleteEntry(params.HTTPRequest.Context(), params.ID)
	if err != nil {
		return entries.NewDeleteEntryNotFound()
	}
	return entries.NewDeleteEntryNoContent()
}

func (h *Handler) ListLinks(params links.ListLinksParams) middleware.Responder {
	result, err := h.svc.ListLinks(params.HTTPRequest.Context())
	if err != nil {
		log.Printf("error listing links: %v", err)
		return links.NewListLinksOK()
	}
	return links.NewListLinksOK().WithPayload(result)
}

func (h *Handler) CreateLink(params links.CreateLinkParams) middleware.Responder {
	result, err := h.svc.CreateLink(params.HTTPRequest.Context(), *params.Body.SourceID, *params.Body.TargetID)
	if err != nil {
		log.Printf("error creating link: %v", err)
		return links.NewCreateLinkCreated()
	}
	return links.NewCreateLinkCreated().WithPayload(result)
}

func (h *Handler) GetLink(params links.GetLinkParams) middleware.Responder {
	result, err := h.svc.GetLink(params.HTTPRequest.Context(), params.ID)
	if err != nil {
		return links.NewGetLinkNotFound()
	}
	return links.NewGetLinkOK().WithPayload(result)
}

func (h *Handler) UpdateLink(params links.UpdateLinkParams) middleware.Responder {
	result, err := h.svc.UpdateLink(params.HTTPRequest.Context(), params.ID, *params.Body.SourceID, *params.Body.TargetID)
	if err != nil {
		return links.NewUpdateLinkNotFound()
	}
	return links.NewUpdateLinkOK().WithPayload(result)
}

func (h *Handler) DeleteLink(params links.DeleteLinkParams) middleware.Responder {
	err := h.svc.DeleteLink(params.HTTPRequest.Context(), params.ID)
	if err != nil {
		return links.NewDeleteLinkNotFound()
	}
	return links.NewDeleteLinkNoContent()
}