package handlers

import (
	"fmt"
	"html"
	"net/http"

	"github.com/joshthewhite/poolvibes/internal/application/command"
	"github.com/joshthewhite/poolvibes/internal/application/services"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	h.renderAuthPage(w, "Sign In", "/login", "Sign In", "/signup", "Don't have an account? Sign up", "")
}

func (h *AuthHandler) SignupPage(w http.ResponseWriter, r *http.Request) {
	h.renderAuthPage(w, "Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", "")
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderAuthPage(w, "Sign In", "/login", "Sign In", "/signup", "Don't have an account? Sign up", "Invalid form data")
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, session, err := h.svc.SignIn(r.Context(), command.SignIn{
		Email:    email,
		Password: password,
	})
	if err != nil {
		h.renderAuthPage(w, "Sign In", "/login", "Sign In", "/signup", "Don't have an account? Sign up", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.renderAuthPage(w, "Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", "Invalid form data")
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirm := r.FormValue("confirm")

	if password != confirm {
		h.renderAuthPage(w, "Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", "Passwords do not match")
		return
	}

	_, session, err := h.svc.SignUp(r.Context(), command.SignUp{
		Email:    email,
		Password: password,
	})
	if err != nil {
		h.renderAuthPage(w, "Sign Up", "/signup", "Sign Up", "/login", "Already have an account? Sign in", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID.String(),
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		_ = h.svc.SignOut(r.Context(), cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) renderAuthPage(w http.ResponseWriter, title, action, buttonText, altURL, altText, errMsg string) {
	errorHTML := ""
	if errMsg != "" {
		errorHTML = fmt.Sprintf(`<div class="notification is-danger is-light">%s</div>`, html.EscapeString(errMsg))
	}

	confirmField := ""
	if title == "Sign Up" {
		confirmField = `
				<div class="field">
					<label class="label">Confirm Password</label>
					<div class="control">
						<input name="confirm" type="password" class="input" placeholder="Confirm your password" required>
					</div>
				</div>`
	}

	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PoolVibes - %s</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@1.0.4/css/bulma.min.css">
</head>
<body>
    <section class="section">
        <div class="container">
            <div class="columns is-centered">
                <div class="column is-4">
                    <div class="box">
                        <h1 class="title has-text-centered">PoolVibes</h1>
                        <h2 class="subtitle has-text-centered">%s</h2>
                        %s
                        <form method="POST" action="%s">
                            <div class="field">
                                <label class="label">Email</label>
                                <div class="control">
                                    <input name="email" type="email" class="input" placeholder="you@example.com" required>
                                </div>
                            </div>
                            <div class="field">
                                <label class="label">Password</label>
                                <div class="control">
                                    <input name="password" type="password" class="input" placeholder="Your password" required>
                                </div>
                            </div>%s
                            <div class="field">
                                <div class="control">
                                    <button type="submit" class="button is-link is-fullwidth">%s</button>
                                </div>
                            </div>
                        </form>
                        <p class="has-text-centered mt-4">
                            <a href="%s">%s</a>
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </section>
</body>
</html>`,
		html.EscapeString(title),
		html.EscapeString(title),
		errorHTML,
		html.EscapeString(action),
		confirmField,
		html.EscapeString(buttonText),
		html.EscapeString(altURL),
		html.EscapeString(altText),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(page))
}
