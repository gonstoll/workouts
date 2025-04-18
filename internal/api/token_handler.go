package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gonstoll/workouts/internal/store"
	"github.com/gonstoll/workouts/internal/tokens"
	"github.com/gonstoll/workouts/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("[ERROR] Decoding on HandleCreateToken: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request"})
		return
	}

	// Get the user and match passwords
	user, err := th.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		th.logger.Printf("[ERROR] GetUserByUsername: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}

	passwordsMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		th.logger.Printf("[ERROR] PasswordHash.Matches: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if !passwordsMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}

	token, err := th.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		th.logger.Printf("[ERROR] CreateNewToken: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
