package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"net/http"

	"github.com/beyachad-maan/auth-service/pkg/api"
	"github.com/beyachad-maan/auth-service/pkg/dao"
	"github.com/beyachad-maan/auth-service/pkg/jwt"
	"github.com/beyachad-maan/auth-service/pkg/password"
)

type LoginHandler struct {
	usersDAO   dao.Users
	privateKey *rsa.PrivateKey
}

func NewLoginHandler(usersDAO dao.Users, privateKey *rsa.PrivateKey) *LoginHandler {
	return &LoginHandler{usersDAO: usersDAO, privateKey: privateKey}
}

func (h *LoginHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var login api.Login
	err = json.Unmarshal(body, &login)
	if err != nil {
		slog.Error("Failed to unmarshal login request", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Authenticate user
	user, err := h.usersDAO.GetUserByUserName(r.Context(), login.Username)
	if err != nil {
		if err == dao.ErrUserNotFound {
			slog.Info("User not found", "user", login.Username)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		slog.Error("Failed to find user", "user", login.Username, "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if !password.VerifyPassword(login.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create JWT token
	token, err := jwt.CreateToken(user.ID.String(), login.Username, h.privateKey)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create JWT token: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	slog.Info(fmt.Sprintf("User '%s' authenticated successfully", login.Username))
	w.WriteHeader(http.StatusOK)
}
