# Chemistry Table Mobile Improvements

## Problem

The chemistry table on mobile has two issues:
1. Action buttons (Plan, Edit, Delete) take too much horizontal space inline
2. Date format ("Jan 2, 3:04 PM") is verbose for a narrow column

## Changes

### 1. Kebab Menu for Actions

Replace the inline Plan/Edit/Delete buttons with a single kebab button (`⋮`) on mobile. The expand chevron stays visible.

**Implementation:**
- Add a Bulma `dropdown is-right` component with class `pv-kebab-menu`
- Toggle via Datastar signal `_chemMenuIdx` (same pattern as `_chemExpandIdx`)
- Only one menu open at a time; clicking another row closes the previous
- Menu items: Plan, Edit, then a divider, then Delete in red
- Desktop: kebab hidden, inline buttons shown (unchanged)
- Mobile: kebab shown, inline buttons hidden
- Uses existing `pv-hidden-mobile` / media query pattern

### 2. Relative Dates

Replace the absolute date format with relative time strings, on both mobile and desktop.

**`relativeTime(t time.Time)` logic:**

| Condition | Output |
|---|---|
| < 1 minute ago | "just now" |
| < 60 minutes ago | "5m ago" |
| < 24 hours ago | "3h ago" |
| Yesterday (calendar day) | "yesterday" |
| 2-6 days ago | "3 days ago" |
| 7-13 days ago | "1 week ago" |
| Same year, older | "Feb 3" |
| Different year | "Feb 3, 2025" |

- Full date/time in `title` attribute for hover tooltip
- Table-driven test for boundary cases
- Added to `helpers.go` alongside existing `dueInText()`

## Files to Modify

- `internal/interface/web/templates/helpers.go` — add `relativeTime()` function
- `internal/interface/web/templates/helpers_test.go` — table-driven tests for `relativeTime()`
- `internal/interface/web/templates/chemistry.templ` — kebab menu markup, relative date in row
- `internal/interface/web/templates/layout.templ` — CSS for `.pv-kebab-menu` visibility

## Approach

- Bulma dropdown + Datastar signal (no new dependencies)
- Server-side relative date computation (no JS libraries)
- Consistent format on mobile and desktop
