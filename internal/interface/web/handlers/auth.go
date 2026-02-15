package handlers

import (
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
	"github.com/joshthewhite/poolvibes/internal/interface/web/templates"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	token := ensureCSRFToken(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.AuthPage("Sign In", "/login", "Sign In", "/signup", "Don't have an account? Sign up", "", token, false).Render(r.Context(), w)
}

func (h *AuthHandler) SignupPage(w http.ResponseWriter, r *http.Request) {
	token := ensureCSRFToken(w, r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.AuthPage("Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", "", token, true).Render(r.Context(), w)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		token := ensureCSRFToken(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.AuthPage("Sign In", "/login", "Sign In", "/signup", "Don't have an account? Sign up", "Invalid form data", token, false).Render(r.Context(), w)
		return
	}

	if !validateCSRF(r) {
		token := ensureCSRFToken(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.AuthPage("Sign In", "/login", "Sign In", "/signup", "Don't have an account? Sign up", "Invalid request, please try again", token, false).Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	_, session, err := h.svc.SignIn(r.Context(), command.SignIn{
		Email:    email,
		Password: password,
	})
	if err != nil {
		token := ensureCSRFToken(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.AuthPage("Sign In", "/login", "Sign In", "/signup", "Don't have an account? Sign up", "Invalid email or password", token, false).Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		token := ensureCSRFToken(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.AuthPage("Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", "Invalid form data", token, true).Render(r.Context(), w)
		return
	}

	if !validateCSRF(r) {
		token := ensureCSRFToken(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.AuthPage("Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", "Invalid request, please try again", token, true).Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if password != confirm {
		token := ensureCSRFToken(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.AuthPage("Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", "Passwords do not match", token, true).Render(r.Context(), w)
		return
	}

	_, session, err := h.svc.SignUp(r.Context(), command.SignUp{
		Email:    email,
		Password: password,
	})
	if err != nil {
		token := ensureCSRFToken(w, r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.AuthPage("Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", err.Error(), token, true).Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		_ = h.svc.SignOut(r.Context(), cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
