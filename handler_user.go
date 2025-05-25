package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iiharsha/rss-go/internal/auth"
	"github.com/iiharsha/rss-go/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("couldn't create user: %v", err))
	}

	responseWithJSON(w, 200, databaseUserToUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apikey)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("couldn't get user: %v", err))
		return
	}

	responseWithJSON(w, 200, databaseUserToUser(user))
}
