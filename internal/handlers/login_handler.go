package handlers

import (
	"log"
	"net/http"
	"web-app/internal/middleware"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := middleware.OAuth2Config.AuthCodeURL("state")

	log.Println("Client ID from oauthconfig: ", middleware.OAuth2Config.ClientID)
	log.Println("Generated URL: ", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
