package handlers

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/beyachad-maan/auth-service/pkg/api"
	"github.com/beyachad-maan/auth-service/pkg/dao"
	"github.com/beyachad-maan/auth-service/pkg/mappers/inbound"
	"github.com/beyachad-maan/auth-service/pkg/mappers/outbound"
	"github.com/beyachad-maan/auth-service/pkg/models"
	"github.com/beyachad-maan/auth-service/pkg/password"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	usersDAO     dao.Users
	jwtPublicKey *rsa.PublicKey
}

func NewUserHandler(usersDAO dao.Users, jwtPublicKey *rsa.PublicKey) *UserHandler {
	return &UserHandler{usersDAO: usersDAO, jwtPublicKey: jwtPublicKey}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Create user
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var user api.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to unmarshal user: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create user object and persist to DB.
	res := inbound.MapUser(user)
	res.ID = models.NewID()
	res.CreatedAt = time.Now()
	passwd, err := password.HashPassword(res.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Password = passwd
	err = h.usersDAO.CreateUser(r.Context(), res)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to create user: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return user object persisted.
	b, err := json.Marshal(outbound.MapUser(res))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to marshal user: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to write user: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	slog.Info("Successfully created user", "username", res.Username, "id", res.ID)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get user
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.usersDAO.GetUserByID(r.Context(), id)
	if err != nil {
		slog.Error("Failed to retrieve user", "error", err)
		if errors.Is(err, dao.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	slog.Info("Successfully retrieved user", "id", user.ID)
	// Return user object persisted.
	b, err := json.Marshal(outbound.MapUser(*user))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	// Get user
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := h.usersDAO.DeleteUserByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
