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
