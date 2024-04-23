package sanitize

import "github.com/microcosm-cc/bluemonday"

// sanitizeString sanitizes the provided string using the UGCPolicy from the bluemonday package.
func SanitizeString(str string) string {
	p := bluemonday.UGCPolicy()

	return p.Sanitize(str)
}