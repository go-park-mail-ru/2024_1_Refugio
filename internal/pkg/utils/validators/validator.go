package validators

import (
	"regexp"
	"strings"
)

// IsEmpty checks if the given string is empty after trimming leading and trailing whitespace.
// Returns true if the string is empty, and false otherwise.
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsValidEmailFormat checks if the provided email string matches the specific format for emails ending with "@mailhub.ru".
func IsValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@mailhub\.su$`)

	return emailRegex.MatchString(email)
}

// IsValidEmailFormatGmail checks if the provided email string matches the specific format for emails ending with "@gmail.com".
func IsValidEmailFormatGmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@gmail\.com$`)

	return emailRegex.MatchString(email)
}
