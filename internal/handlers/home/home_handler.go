package home

import (
	"net/http"
	"web-app/internal/middleware"
	"web-app/internal/services"
	"web-app/pkg/utils"
)

type HomeHandler struct {
	sessionService *services.SessionService
}

func NewHomeHandler(sessionService *services.SessionService) *HomeHandler {
	return &HomeHandler{sessionService: sessionService}
}

func (h *HomeHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.FromContext(r.Context())

	pageData := struct{}{}

	tmplData := utils.TemplateData{
		BuildTime:  utils.BuildTime,
		IsLoggedIn: ok,
		ShowHeader: true,
		Page:       "home",
		Data:       pageData,
	}

	utils.RenderTemplate(w, "home.html", tmplData)
}
