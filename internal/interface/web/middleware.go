package web

import (
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/services"
)

func requireAuth(authSvc *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user, err := authSvc.GetUserBySession(r.Context(), cookie.Value)
		if err != nil || user == nil {
			http.SetCookie(w, &http.Cookie{
				Name:   "session_id",
				Value:  "",
				Path:   "/",
				MaxAge: -1,
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := services.WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func requireAdmin(authSvc *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return requireAuth(authSvc, func(w http.ResponseWriter, r *http.Request) {
		user, err := services.UserFromContext(r.Context())
		if err != nil || !user.IsAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
