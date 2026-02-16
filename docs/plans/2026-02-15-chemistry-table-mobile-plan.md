# Chemistry Table Mobile Improvements — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Improve the chemistry table mobile UX by replacing inline action buttons with a kebab menu and switching to relative date formatting.

**Architecture:** Server-side relative date helper in Go, Bulma dropdown toggled via Datastar signal for the kebab menu. No new dependencies. CSS media queries control mobile/desktop visibility.

**Tech Stack:** Go, templ, Bulma CSS, Datastar signals

---

### Task 1: Add `relativeTime()` helper — write failing tests

**Files:**
- Create: `internal/interface/web/templates/helpers_test.go`

**Step 1: Write the failing tests**

Create `internal/interface/web/templates/helpers_test.go` with table-driven tests. Use a `now` parameter so tests are deterministic (the function will accept `now` as a second arg).

```go
package templates

import (
	"testing"
	"time"
)

func TestRelativeTime(t *testing.T) {
	// Fixed reference time: Feb 15, 2026 14:00:00 UTC
	now := time.Date(2026, 2, 15, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{"just now", now.Add(-30 * time.Second), "just now"},
		{"1 minute ago", now.Add(-1 * time.Minute), "1m ago"},
		{"5 minutes ago", now.Add(-5 * time.Minute), "5m ago"},
		{"59 minutes ago", now.Add(-59 * time.Minute), "59m ago"},
		{"1 hour ago", now.Add(-1 * time.Hour), "1h ago"},
		{"3 hours ago", now.Add(-3 * time.Hour), "3h ago"},
		{"23 hours ago", now.Add(-23 * time.Hour), "23h ago"},
		{"yesterday", time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC), "yesterday"},
		{"2 days ago", time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC), "2 days ago"},
		{"6 days ago", time.Date(2026, 2, 9, 10, 0, 0, 0, time.UTC), "6 days ago"},
		{"1 week ago", time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC), "1 week ago"},
		{"13 days ago", time.Date(2026, 2, 2, 10, 0, 0, 0, time.UTC), "1 week ago"},
		{"same year older", time.Date(2026, 1, 10, 10, 0, 0, 0, time.UTC), "Jan 10"},
		{"different year", time.Date(2025, 6, 15, 10, 0, 0, 0, time.UTC), "Jun 15, 2025"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := relativeTimeFrom(tt.t, now); got != tt.want {
				t.Errorf("relativeTimeFrom() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test -v -run TestRelativeTime ./internal/interface/web/templates/`
Expected: FAIL — `relativeTimeFrom` is undefined

---

### Task 2: Implement `relativeTime()` helper

**Files:**
- Modify: `internal/interface/web/templates/helpers.go`

**Step 1: Add `relativeTimeFrom()` and `relativeTime()` to helpers.go**

Add after the existing `dueInText()` function:

```go
// relativeTime returns a human-friendly relative time string for a past timestamp.
func relativeTime(t time.Time) string {
	return relativeTimeFrom(t, time.Now())
}

// relativeTimeFrom returns a relative time string using the given reference time.
// Exported for testing.
func relativeTimeFrom(t, now time.Time) string {
	d := now.Sub(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}

	// Calendar-day logic
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	days := int(today.Sub(tDay).Hours() / 24)

	switch {
	case days == 1:
		return "yesterday"
	case days <= 6:
		return fmt.Sprintf("%d days ago", days)
	case days <= 13:
		return "1 week ago"
	case t.Year() == now.Year():
		return t.Format("Jan 2")
	default:
		return t.Format("Jan 2, 2006")
	}
}
```

**Step 2: Run tests to verify they pass**

Run: `go test -v -run TestRelativeTime ./internal/interface/web/templates/`
Expected: PASS — all cases green

**Step 3: Commit**

```bash
git add internal/interface/web/templates/helpers.go internal/interface/web/templates/helpers_test.go
git commit -m "feat: add relativeTime helper for human-friendly timestamps"
```

---

### Task 3: Add kebab menu CSS

**Files:**
- Modify: `internal/interface/web/templates/layout.templ` (CSS section around line 502-522)

**Step 1: Add kebab menu visibility rules**

In the `@media screen and (max-width: 768px)` block (around line 503), add:

```css
.pv-kebab-menu {
    display: inline-flex !important;
}
.pv-action-btn-desktop {
    display: none !important;
}
```

In the `@media screen and (min-width: 769px)` block (around line 515), add:

```css
.pv-kebab-menu {
    display: none !important;
}
```

Also add (outside any media query, before the closing `</style>`) dropdown positioning fix for table cells:

```css
.pv-kebab-menu .dropdown-menu {
    min-width: 8rem;
}
```

**Step 2: Commit**

```bash
git add internal/interface/web/templates/layout.templ
git commit -m "feat: add CSS for mobile kebab menu visibility"
```

---

### Task 4: Update chemistry table template

**Files:**
- Modify: `internal/interface/web/templates/chemistry.templ`

**Step 1: Add `_chemMenuIdx` signal**

On line 19, after `data-signals:_chemExpandIdx="'none'"`, add:

```
data-signals:_chemMenuIdx="'none'"
```

**Step 2: Update date cell in `chemistryRow`**

Replace line 137:
```
<td>{ l.TestedAt.Format("Jan 2, 3:04 PM") }</td>
```
With:
```
<td title={ l.TestedAt.Format("Jan 2, 2006 3:04 PM") }>{ relativeTime(l.TestedAt) }</td>
```

**Step 3: Replace action buttons with kebab menu + desktop buttons**

Replace lines 145-158 (the entire actions `<td>`) with:

```html
<td class="has-text-right">
    <div class="buttons is-right are-small" style="flex-wrap: nowrap;">
        <button
            data-on:click={ fmt.Sprintf("$_chemExpandIdx = $_chemExpandIdx === '%d' ? 'none' : '%d'", idx, idx) }
            class="button is-small pv-expand-btn"
        >
            <span data-class:is-hidden={ fmt.Sprintf("$_chemExpandIdx === '%d'", idx) }>&#9660;</span>
            <span class="is-hidden" data-class:is-hidden={ fmt.Sprintf("$_chemExpandIdx !== '%d'", idx) }>&#9650;</span>
        </button>
        <!-- Mobile: kebab menu -->
        <div class="dropdown is-right pv-kebab-menu" data-class:is-active={ fmt.Sprintf("$_chemMenuIdx === '%d'", idx) }>
            <div class="dropdown-trigger">
                <button
                    class="button is-small"
                    data-on:click={ fmt.Sprintf("$_chemMenuIdx = $_chemMenuIdx === '%d' ? 'none' : '%d'", idx, idx) }
                    aria-haspopup="true"
                >
                    <span>&#8942;</span>
                </button>
            </div>
            <div class="dropdown-menu" role="menu">
                <div class="dropdown-content">
                    <a class="dropdown-item" data-on:click={ "@get('/chemistry/" + l.ID.String() + "/plan')" }>Plan</a>
                    <a class="dropdown-item" data-on:click={ "@get('/chemistry/" + l.ID.String() + "/edit')" }>Edit</a>
                    <hr class="dropdown-divider"/>
                    <a class="dropdown-item has-text-danger" data-on:click={ "@delete('/chemistry/" + l.ID.String() + "')" }>Delete</a>
                </div>
            </div>
        </div>
        <!-- Desktop: inline buttons -->
        <button data-on:click={ "@get('/chemistry/" + l.ID.String() + "/plan')" } class="button is-info is-outlined is-small pv-action-btn-desktop">Plan</button>
        <button data-on:click={ "@get('/chemistry/" + l.ID.String() + "/edit')" } class="button is-primary is-outlined is-small pv-action-btn-desktop">Edit</button>
        <button data-on:click={ "@delete('/chemistry/" + l.ID.String() + "')" } class="button is-danger is-outlined is-small pv-action-btn-desktop">Delete</button>
    </div>
</td>
```

**Step 4: Generate templ**

Run: `task templ`

**Step 5: Run all tests**

Run: `task test`
Expected: all pass

**Step 6: Commit**

```bash
git add internal/interface/web/templates/chemistry.templ internal/interface/web/templates/chemistry_templ.go
git commit -m "feat: add kebab menu and relative dates to chemistry table"
```

---

### Task 5: Manual verification

**Step 1: Start dev server**

Run: `task dev`

**Step 2: Verify desktop**

- Chemistry table shows relative dates ("3h ago", "yesterday", "Feb 3", etc.)
- Hover over date shows full date/time tooltip
- Inline Plan/Edit/Delete buttons visible, no kebab button
- Expand chevron is hidden

**Step 3: Verify mobile**

- Use browser dev tools to switch to mobile viewport (375px)
- Only Date, pH, FC, Actions columns visible
- Expand chevron visible, inline Plan/Edit/Delete hidden
- Kebab button (`⋮`) visible
- Tap kebab: dropdown appears with Plan, Edit, divider, Delete (red)
- Tap another row's kebab: previous menu closes, new one opens
- Tap same kebab: menu closes
- Plan/Edit/Delete actions work from the dropdown

**Step 4: Stop dev server**
