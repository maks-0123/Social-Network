package main

import (
	"log"
	"my-social-network/internal/database"
	"my-social-network/internal/handlers"
	"my-social-network/internal/middleware"
	"my-social-network/internal/repository"
	"my-social-network/internal/service"

	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

func main() {
	pool, err := database.NewConnection()
	if err != nil {
		log.Println("Ошибка подключения к БД", err)
	}
	defer pool.Close()
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	userRepo := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, sessionManager)
	postRepo := repository.NewPostRepository(pool)
	postService := service.NewPostService(postRepo)
	postHandler := handlers.NewPostHandler(postService, sessionManager)

	r := chi.NewRouter()

	r.Use(sessionManager.LoadAndSave)

	r.Post("/user/register", userHandler.RegisterHandler)
	r.Post("/user/login", userHandler.LoginHandler)

	r.Post("/posts/create", postHandler.CreatePostHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Главная страница"))
	})
	r.Get("/feed", postHandler.GetFeedHandler)
	r.With(middleware.RequireAuth(sessionManager)).Get("/profile", func(w http.ResponseWriter, r *http.Request) {
		userID := sessionManager.GetInt(r.Context(), "userID")
		w.Write([]byte("Это защищенная страница профиля! UserID: " + string(rune(userID))))
	})
	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
