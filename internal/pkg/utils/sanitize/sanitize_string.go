package sanitize

import "github.com/microcosm-cc/bluemonday"

// SanitizeString sanitizes the provided string using the UGCPolicy.
func SanitizeString(str string) string {
	p := bluemonday.UGCPolicy()

	return p.Sanitize(str)
}
