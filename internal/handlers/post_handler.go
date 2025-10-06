package handlers

import (
	"encoding/json"
	"my-social-network/internal/service"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
)

type PostHandler struct {
	postService    *service.PostService
	sessionManager *scs.SessionManager
}

func NewPostHandler(postService *service.PostService, sessionManager *scs.SessionManager) *PostHandler {
	return &PostHandler{
		postService:    postService,
		sessionManager: sessionManager,
	}
}

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Проверить что метод POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Получить userID из сессии
	userID := h.sessionManager.GetInt(r.Context(), "userID")
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. Получить content из JSON тела
	var request struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 4. Вызвать сервис
	if err := h.postService.CreatePost(userID, request.Content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5. Вернуть успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Post created successfully",
	})
}

func (h *PostHandler) GetFeedHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получить limit из query параметра
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // значение по умолчанию

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 2. Получить посты
	posts, err := h.postService.GetFeed(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Вернуть посты в JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
