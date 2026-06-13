package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prajwal-huggi/backend_go/internal/domain"
	"github.com/prajwal-huggi/backend_go/internal/dto"
	"github.com/prajwal-huggi/backend_go/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// CreateUser godoc
//
// @Summary Create a new user
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body domain.UserModel true "User Payload"
// @Success 201 {object} domain.UserModel
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.UserModel

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.CreateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUsers godoc
//
// @Summary Get all users
// @Description Returns all users
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.UserModel
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /users [get]

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser godoc
//
// @Summary Get user by ID
// @Description Returns a single user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} domain.UserModel
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser godoc
//
// @Summary Update user
// @Description Update an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body domain.UserModel true "User Payload"
// @Success 200 {object} domain.UserModel
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user domain.UserModel

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// DeleteUser godoc
//
// @Summary Delete user
// @Description Delete user by ID
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Login godoc
//
// @Summary Login User
// @Description Authenticate user and return JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginRequest

	json.NewDecoder(r.Body).Decode(&req)

	access, refresh, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

// Refresh godoc
//
// @Summary Refresh access token
// @Description Generate a new access token using a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh Request"
// @Success 200 {object} dto.RefreshResponse
// @Failure 401 {string} string
// @Router /refresh [post]
func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest

	json.NewDecoder(r.Body).Decode(&req)

	access, err := h.service.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid token", 401)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": access,
	})
}
