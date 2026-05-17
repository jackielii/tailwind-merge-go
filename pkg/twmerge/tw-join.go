package twmerge

import "strings"

// ClassNameValue can be supplied to TwJoin / Merge. Supported runtime types:
//   - string
//   - []ClassNameValue
//   - []string
//   - []any (each element handled recursively)
//   - nil / bool false / 0 (skipped)
type ClassNameValue = any

// TwJoin concatenates class-name values, skipping falsy entries.
// Mirrors the lukeed/clsx behaviour the JS package's tw-join inherits from.
func TwJoin(values ...ClassNameValue) string {
	var b strings.Builder
	for _, v := range values {
		if v == nil {
			continue
		}
		s := toClassString(v)
		if s == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(s)
	}
	return b.String()
}

func toClassString(v any) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return val
	case bool:
		return ""
	case int:
		if val == 0 {
			return ""
		}
		// Non-zero ints are explicitly unsupported in JS twJoin (only 0/0n are falsy);
		// in Go we follow the spirit and ignore them.
		return ""
	case []string:
		var b strings.Builder
		for _, s := range val {
			if s == "" {
				continue
			}
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(s)
		}
		return b.String()
	case []any:
		var b strings.Builder
		for _, e := range val {
			s := toClassString(e)
			if s == "" {
				continue
			}
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(s)
		}
		return b.String()
	default:
		return ""
	}
}
