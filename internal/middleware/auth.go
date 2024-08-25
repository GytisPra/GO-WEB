package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"web-app/internal/models"
	"web-app/internal/services"
	"web-app/pkg/utils"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var OAuth2Config *oauth2.Config

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	OAuth2Config = &oauth2.Config{
		ClientID:     os.Getenv("DISCORD_CLIENT_ID"),
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:3000/callback/discord",
		Scopes:       []string{"identify", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}
}

type CallbackHandler struct {
	userSerivce    *services.UserService
	accountService *services.AcountService
	sessionService *services.SessionService
}

func NewCallbackHandler(userSerivce *services.UserService, accountService *services.AcountService, sessionService *services.SessionService) *CallbackHandler {
	return &CallbackHandler{userSerivce: userSerivce, accountService: accountService, sessionService: sessionService}
}

func (h *CallbackHandler) DiscordCallbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	if code == "" {
		utils.HandleError(w, "Missing code", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := OAuth2Config.Exchange(r.Context(), code)
	if err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	// Use the access token to get user information
	client := OAuth2Config.Client(r.Context(), token)
	resp, err := client.Get("https://discord.com/api/v10/users/@me")
	if err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.HandleError(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Read and parse the user info response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
		return
	}

	var discordUser map[string]interface{}
	if err := json.Unmarshal(body, &discordUser); err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to parse user info: %v", err), http.StatusInternalServerError)
		return
	}

	// discordID := discordUser["id"].(string)
	discordID, ok := discordUser["id"].(string)
	if !ok || discordID == "" {
		utils.HandleError(w, "Failed to retrieve Discord ID", http.StatusInternalServerError)
		return
	}

	username, ok := discordUser["username"].(string)
	if !ok || username == "" {
		utils.HandleError(w, "Failed to retrieve username", http.StatusInternalServerError)
		return
	}

	email, ok := discordUser["email"].(string)
	if !ok || email == "" {
		utils.HandleError(w, "Failed to retrieve email", http.StatusInternalServerError)
		return
	}

	var user *models.User
	user, err = h.userSerivce.CreateUser(username, email)
	if err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("User ID: ", user.ID)

	var session *models.Session
	session, err = h.sessionService.CreateSession(user.ID)
	if err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to create session: %v", err), http.StatusInternalServerError)
		return
	}

	scopes, ok := token.Extra("scope").(string)
	if !ok || scopes == "" {
		utils.HandleError(w, "Failed to retrieve 'scopes' from token", http.StatusInternalServerError)
		return
	}

	expiresIn, ok := token.Extra("expires_in").(float64)
	if !ok {
		utils.HandleError(w, "Failed to retrieve 'expires_in' from token", http.StatusInternalServerError)
		return
	}

	tokenType, ok := token.Extra("token_type").(string)
	if !ok || tokenType == "" {
		utils.HandleError(w, "Failed to retrieve 'token_type' from token", http.StatusInternalServerError)
		return
	}

	var account *models.Account
	account, err = h.accountService.CreateAccount(user.ID, "oauth", "discord", discordID, &token.RefreshToken, &token.AccessToken, &expiresIn, &tokenType, &scopes, nil, nil, nil)
	if err != nil {
		utils.HandleError(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":    user,
		"account": account,
		"session": session,
	})
}
