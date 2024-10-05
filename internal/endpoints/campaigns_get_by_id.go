package endpoints

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) CampaignGetById(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	id := chi.URLParam(r, "id")
	campaign, err := h.CampaignService.GetBy(id)
	if err == nil && campaign == nil {
		return nil, http.StatusNotFound, err
	}
	fmt.Println("hendeler " + id)
	return campaign, 200, err
}
