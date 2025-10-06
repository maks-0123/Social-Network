package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func RequireAuth(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sessionManager.GetInt(r.Context(), "userID")
			if userID == 0 {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
