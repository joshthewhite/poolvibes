# Marketing Landing Page Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add a marketing landing page at `/` for unauthenticated visitors that showcases PoolVibes features and drives signups.

**Architecture:** Add an `optionalAuth` middleware that sets the user in context if authenticated but doesn't redirect if not. The root handler checks for a user: authenticated users get the dashboard, unauthenticated users get a standalone landing page template (like `auth.templ` — full HTML document, not using `Layout`).

**Tech Stack:** Go handlers, templ templates, Bulma CSS, inline styles matching the app's design system.

**Design doc:** `docs/plans/2026-02-18-marketing-page-design.md`

---

### Task 1: Add `optionalAuth` middleware

**Files:**
- Modify: `internal/interface/web/middleware.go`

**Step 1: Add `optionalAuth` function**

Add after `requireAdmin` in `middleware.go`:

```go
func optionalAuth(authSvc *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil && cookie.Value != "" {
			user, err := authSvc.GetUserBySession(r.Context(), cookie.Value)
			if err == nil && user != nil {
				ctx := services.WithUser(r.Context(), user)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	}
}
```

**Step 2: Commit**

```bash
git add internal/interface/web/middleware.go
git commit -m "feat: add optionalAuth middleware for landing page"
```

---

### Task 2: Add `Root` method to `PageHandler`

**Files:**
- Modify: `internal/interface/web/handlers/page.go`

**Step 1: Add `Root` method**

Add a `Root` method that checks for an authenticated user and dispatches accordingly:

```go
func (h *PageHandler) Root(w http.ResponseWriter, r *http.Request) {
	user, _ := services.UserFromContext(r.Context())
	if user != nil {
		h.Index(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.LandingPage().Render(r.Context(), w)
}
```

**Step 2: Commit**

```bash
git add internal/interface/web/handlers/page.go
git commit -m "feat: add Root handler to dispatch landing vs dashboard"
```

---

### Task 3: Wire up routing

**Files:**
- Modify: `internal/interface/web/server.go`

**Step 1: Add `maybeAuth` closure and update root route**

In `setupRoutes()`, add the `maybeAuth` closure alongside `auth` and `admin` (line 51):

```go
auth      := func(h http.HandlerFunc) http.HandlerFunc { return requireAuth(s.authSvc, h) }
admin     := func(h http.HandlerFunc) http.HandlerFunc { return requireAdmin(s.authSvc, h) }
maybeAuth := func(h http.HandlerFunc) http.HandlerFunc { return optionalAuth(s.authSvc, h) }
```

Change the root route (line 65) from:
```go
s.mux.HandleFunc("GET /{$}", auth(pageHandler.Index))
```
to:
```go
s.mux.HandleFunc("GET /{$}", maybeAuth(pageHandler.Root))
```

**Step 2: Commit**

```bash
git add internal/interface/web/server.go
git commit -m "feat: route unauthenticated root to landing page"
```

---

### Task 4: Create landing page template

**Files:**
- Create: `internal/interface/web/templates/landing.templ`

> **For Claude:** REQUIRED SUB-SKILL: Use frontend-design to create the landing page template. Reference the design doc at `docs/plans/2026-02-18-marketing-page-design.md` for the full visual spec.

This is a standalone HTML document (like `auth.templ`, NOT using `Layout`). It must include:

1. **HTML head**: Bulma CDN, Inter/Inter Tight fonts, viewport meta, inline `<style>` block with design tokens from the app
2. **Landing navbar**: Logo/brand left, "Log In" + "Sign Up" buttons right
3. **Hero section**: Dark bg (`#13111C`), bold headline, subline, teal "Get Started" CTA → `/signup`, secondary "Log in" link, wave/gradient transition to light body
4. **Feature sections** (5 alternating left/right):
   - Water Chemistry: "Know your water like never before"
   - Maintenance Tasks: "Never forget to skim again"
   - Chemical Inventory: "Running low? You'll know first."
   - Equipment Tracking: "Your pump has a story. Keep it."
   - Smart Notifications: "A nudge before things go sideways"
   - Each with CSS-rendered visual mockups using app design tokens (neumorphic cards, status colors)
5. **Stats bar**: Light teal bg, 3 stats: "6 key parameters", "Zero guesswork", "Free & open source"
6. **Final CTA**: Dark bg, "Your pool deserves better than a spreadsheet.", "Get Started" button
7. **Footer**: Minimal, brand name, "Free & open source" note
8. **Responsive**: Full mobile support, stacked layouts on small screens
9. **Dark mode**: `prefers-color-scheme` media query support

The templ component signature:

```go
templ LandingPage() {
    <!DOCTYPE html>
    <html lang="en">
    // ... full standalone document
    </html>
}
```

**Step 2: Generate templ code**

Run: `task templ`

**Step 3: Verify it compiles**

Run: `task build`

**Step 4: Commit**

```bash
git add internal/interface/web/templates/landing.templ internal/interface/web/templates/landing_templ.go
git commit -m "feat: add marketing landing page template"
```

---

### Task 5: Visual QA and polish

**Step 1: Run the dev server**

Run: `task dev`

**Step 2: Test unauthenticated flow**

- Open `http://localhost:8080/` in a browser (not logged in)
- Verify the landing page renders with all sections
- Verify "Get Started" links to `/signup`
- Verify "Log In" links to `/login`
- Test responsive layout (resize browser)
- Test dark mode (toggle OS preference)

**Step 3: Test authenticated flow**

- Log in at `/login`
- Navigate to `/` — verify you see the dashboard, not the landing page

**Step 4: Fix any visual issues found during QA**

Run `task templ && task build` after any template changes.

**Step 5: Commit any fixes**

```bash
git add internal/interface/web/templates/landing.templ internal/interface/web/templates/landing_templ.go
git commit -m "fix: landing page visual polish"
```

---

### Task 6: Update auth page redirects

**Files:**
- Modify: `internal/interface/web/middleware.go`

Currently, `requireAuth` redirects unauthenticated users to `/login`. This is fine — authenticated-only routes should still redirect to login, not the landing page. No change needed.

However, the login and signup pages link back to each other but have no link to the landing page. Consider adding a "Back to home" or brand logo link on the auth pages that links to `/`.

**Step 1: Update auth template (optional)**

If the auth template (`auth.templ`) brand name at the top is not already a link, wrap it in `<a href="/">`:

```go
<a href="/" class="pv-auth-brand">PoolVibes</a>
```

**Step 2: Generate and commit**

```bash
task templ
git add internal/interface/web/templates/auth.templ internal/interface/web/templates/auth_templ.go
git commit -m "feat: link auth page brand to landing page"
```

---

### Task 7: Final verification

**Step 1: Run full test suite**

Run: `task test`
Expected: All tests pass

**Step 2: Run linter**

Run: `task lint`
Expected: No errors

**Step 3: Full flow test**

1. Start fresh: `task dev`
2. Open `/` unauthenticated → see landing page
3. Click "Get Started" → go to `/signup`
4. Sign up → go to dashboard
5. Visit `/` while logged in → see dashboard
6. Log out → redirected to `/login`
7. Visit `/` → see landing page

**Step 4: Final commit if needed**

```bash
git add -A
git commit -m "chore: landing page final cleanup"
```
