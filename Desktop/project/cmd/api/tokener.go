package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type TokenRequest struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scope        []string `json:"scope"`
}

type TokenResponse struct {
	AccessToken   string   `json:"access_token"`
	ExpiresIn     int      `json:"expires_in"`
	RefreshToken  string   `json:"refresh_token"`
	Scope         []string `json:"scope"`
	SecurityLevel string   `json:"security_level"`
	TokenType     string   `json:"token_type"`
}

func (app *application) TokenGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON input
	var req TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate client credentials
	var dbClientSecret string
	err := app.DB.QueryRow(
		"SELECT client_secret FROM public.user WHERE client_id = $1",
		req.ClientID,
	).Scan(&dbClientSecret)

	if err == sql.ErrNoRows || dbClientSecret != req.ClientSecret {
		http.Error(w, "Invalid client credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Generate new access token and expiration time
	accessToken := generateAccessToken()
	expirationTime := time.Now().Add(2 * time.Hour)

	// Insert the token into the database
	_, err = app.DB.Exec(
		`INSERT INTO public.token (client_id, access_scope, access_token, expiration_time)
		 VALUES ($1, $2, $3, $4)`,
		req.ClientID, req.Scope, accessToken, expirationTime,
	)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	resp := TokenResponse{
		AccessToken:   accessToken,
		ExpiresIn:     7200,
		RefreshToken:  "",
		Scope:         req.Scope,
		SecurityLevel: "normal",
		TokenType:     "Bearer",
	}

	// Set response headers and send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func generateAccessToken() string {
	// Generate a 22-character random token using PostgreSQL substring logic
	var token string
	err := app.DB.QueryRow(
		`SELECT SUBSTR(UPPER(md5(random()::text)), 2, 22)`,
	).Scan(&token)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate access token: %v", err))
	}
	return token
}
