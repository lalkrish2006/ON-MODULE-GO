package middleware

import (
	"net/http"
	"od-system/internal/services"
)

// AuthMiddleware checks if the user is logged in
func RequireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := services.GetSession(r)
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware checks if the user has the required role
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := services.GetSession(r)
		userRole, ok := session.Values["role"].(string)
		if !ok || userRole != role {
			http.Redirect(w, r, "/login?error=unauthorized", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
