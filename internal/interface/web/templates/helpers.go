package templates

import (
	"fmt"
	"strings"
	"time"
)

func escapeJS(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}

func valueClass(inRange bool) string {
	if inRange {
		return "value-ok"
	}
	return "value-warn"
}

func fmtDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func fmtFloat(f float64, prec int) string {
	return fmt.Sprintf("%.*f", prec, f)
}

func fmtFloatG(f float64) string {
	return fmt.Sprintf("%g", f)
}

func fmtGallons(n int) string {
	s := fmt.Sprintf("%d", n)
	if n < 1000 {
		return s
	}
	// Insert commas from right to left
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

func dueInText(dueDate time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	due := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())

	days := int(due.Sub(today).Hours() / 24)

	switch {
	case days == 0:
		return "due today"
	case days == 1:
		return "due tomorrow"
	case days > 1 && days <= 13:
		return fmt.Sprintf("due in %d days", days)
	case days > 13 && days <= 8*7:
		weeks := days / 7
		if weeks == 1 {
			return "due in 1 week"
		}
		return fmt.Sprintf("due in %d weeks", weeks)
	case days > 8*7:
		months := days / 30
		if months <= 1 {
			return "due in 1 month"
		}
		return fmt.Sprintf("due in %d months", months)
	case days == -1:
		return "overdue 1 day"
	case days < -1 && days >= -13:
		return fmt.Sprintf("overdue %d days", -days)
	case days < -13 && days >= -8*7:
		weeks := -days / 7
		if weeks == 1 {
			return "overdue 1 week"
		}
		return fmt.Sprintf("overdue %d weeks", weeks)
	default:
		months := -days / 30
		if months <= 1 {
			return "overdue 1 month"
		}
		return fmt.Sprintf("overdue %d months", months)
	}
}

// relativeTime returns a human-friendly relative time string for a past timestamp.
func relativeTime(t time.Time) string {
	return relativeTimeFrom(t, time.Now())
}

// relativeTimeFrom returns a relative time string using the given reference time.
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

func statusColor(status string) string {
	switch status {
	case "danger":
		return "has-text-danger"
	case "warning":
		return "has-text-warning"
	default:
		return "has-text-success"
	}
}

func demoExpiryText(expiresAt *time.Time) string {
	if expiresAt == nil {
		return ""
	}
	remaining := time.Until(*expiresAt)
	if remaining <= 0 {
		return "expired"
	}
	hours := int(remaining.Hours())
	if hours >= 1 {
		return fmt.Sprintf("expires in %dh", hours)
	}
	minutes := int(remaining.Minutes())
	return fmt.Sprintf("expires in %dm", minutes)
}

// sortAction returns the Datastar action string for a sortable column header click.
// It toggles direction if already sorting by this column, otherwise sorts desc.
func sortAction(col, currentSortBy, currentSortDir string) string {
	newDir := "desc"
	if col == currentSortBy && currentSortDir == "desc" {
		newDir = "asc"
	}
	return fmt.Sprintf("$chemsortby='%s'; $chemsortdir='%s'; $chempage=1; @get('/chemistry')", col, newDir)
}

// sortIndicator returns an arrow character for the active sort column.
func sortIndicator(col, currentSortBy, currentSortDir string) string {
	if col != currentSortBy {
		return ""
	}
	if currentSortDir == "asc" {
		return " \u2191"
	}
	return " \u2193"
}

// paginationPages returns page numbers to display with ellipsis gaps.
// Returns numbers 1-based; 0 represents an ellipsis.
func paginationPages(current, total int) []int {
	if total <= 7 {
		pages := make([]int, total)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}
	var pages []int
	pages = append(pages, 1)
	if current > 3 {
		pages = append(pages, 0) // ellipsis
	}
	for p := current - 1; p <= current+1; p++ {
		if p > 1 && p < total {
			pages = append(pages, p)
		}
	}
	if current < total-2 {
		pages = append(pages, 0) // ellipsis
	}
	pages = append(pages, total)
	return pages
}

// showingRange returns "Showing X-Y of Z" text.
func showingRange(page, pageSize, totalItems int) string {
	if totalItems == 0 {
		return "No results"
	}
	start := (page-1)*pageSize + 1
	end := start + pageSize - 1
	if end > totalItems {
		end = totalItems
	}
	return fmt.Sprintf("Showing %d\u2013%d of %d", start, end, totalItems)
}

func dueInClass(dueDate time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	due := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())

	days := int(due.Sub(today).Hours() / 24)

	switch {
	case days < 0:
		return "tag is-danger is-light"
	case days <= 1:
		return "tag is-warning"
	case days <= 3:
		return "tag is-warning is-light"
	default:
		return "tag is-light"
	}
}
