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
		return "has-text-success has-text-weight-semibold"
	}
	return "has-text-danger has-text-weight-bold"
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
