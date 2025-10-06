package handlers

import (
	"encoding/json"
	"my-social-network/internal/service"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type UserHandler struct {
	userService    *service.UserService
	sessionManager *scs.SessionManager
}

func NewUserHandler(userService *service.UserService, sessionManager *scs.SessionManager) *UserHandler {
	return &UserHandler{
		userService:    userService,
		sessionManager: sessionManager,
	}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	err := h.userService.RegisterUser(request.Email, request.Username, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Пользователь успешно зарегистрирован",
	})
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	user, err := h.userService.LoginUser(request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Сохраняем userID в сессии
	h.sessionManager.Put(r.Context(), "userID", user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Вход выполнен успешно",
		"user":    user,
	})
}
