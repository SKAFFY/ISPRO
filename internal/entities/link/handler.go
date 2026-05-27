package link

import (
	"log"

	"github.com/go-openapi/runtime/middleware"

	"github.com/mkheyfets/ispro-app/internal/entities/link/domain"
	"github.com/mkheyfets/ispro-app/internal/models"
	"github.com/mkheyfets/ispro-app/internal/restapi/operations/links"
)

type Handler struct {
	listLinks  *ListLinksUseCase
	getLink    *GetLinkUseCase
	createLink *CreateLinkUseCase
	updateLink *UpdateLinkUseCase
	deleteLink *DeleteLinkUseCase
}

func NewHandler(
	listLinks *ListLinksUseCase,
	getLink *GetLinkUseCase,
	createLink *CreateLinkUseCase,
	updateLink *UpdateLinkUseCase,
	deleteLink *DeleteLinkUseCase,
) *Handler {
	return &Handler{
		listLinks:  listLinks,
		getLink:    getLink,
		createLink: createLink,
		updateLink: updateLink,
		deleteLink: deleteLink,
	}
}

func (h *Handler) ListLinks(params links.ListLinksParams) middleware.Responder {
	result, err := h.listLinks.Handle(params.HTTPRequest.Context())
	if err != nil {
		log.Printf("error listing links: %v", err)
		return links.NewListLinksOK()
	}
	return links.NewListLinksOK().WithPayload(toModelList(result))
}

func (h *Handler) CreateLink(params links.CreateLinkParams) middleware.Responder {
	cmd := CreateLinkCommand{
		SourceID: *params.Body.SourceID,
		TargetID: *params.Body.TargetID,
	}
	result, err := h.createLink.Handle(params.HTTPRequest.Context(), cmd)
	if err != nil {
		log.Printf("error creating link: %v", err)
		return links.NewCreateLinkCreated()
	}
	return links.NewCreateLinkCreated().WithPayload(toModel(result))
}

func (h *Handler) GetLink(params links.GetLinkParams) middleware.Responder {
	query := GetLinkQuery{ID: params.ID}
	result, err := h.getLink.Handle(params.HTTPRequest.Context(), query)
	if err != nil {
		return links.NewGetLinkNotFound()
	}
	return links.NewGetLinkOK().WithPayload(toModel(result))
}

func (h *Handler) UpdateLink(params links.UpdateLinkParams) middleware.Responder {
	cmd := UpdateLinkCommand{
		ID:       params.ID,
		SourceID: *params.Body.SourceID,
		TargetID: *params.Body.TargetID,
	}
	result, err := h.updateLink.Handle(params.HTTPRequest.Context(), cmd)
	if err != nil {
		return links.NewUpdateLinkNotFound()
	}
	return links.NewUpdateLinkOK().WithPayload(toModel(result))
}

func (h *Handler) DeleteLink(params links.DeleteLinkParams) middleware.Responder {
	cmd := DeleteLinkCommand{ID: params.ID}
	err := h.deleteLink.Handle(params.HTTPRequest.Context(), cmd)
	if err != nil {
		return links.NewDeleteLinkNotFound()
	}
	return links.NewDeleteLinkNoContent()
}

func toModel(l *domain.Link) *models.Link {
	if l == nil {
		return nil
	}
	return &models.Link{
		ID:       l.ID,
		SourceID: l.SourceID,
		TargetID: l.TargetID,
	}
}

func toModelList(links []*domain.Link) []*models.Link {
	result := make([]*models.Link, len(links))
	for i, l := range links {
		result[i] = toModel(l)
	}
	return result
}
