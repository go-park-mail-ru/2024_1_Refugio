package models

type UserGender string

const (
	Male   UserGender = "Male"
	Female UserGender = "Female"
	Other  UserGender = "Other"
)

// Function to check if the given value is a valid UserGender.
func IsValidGender(gender UserGender) bool {
	switch gender {
	case Male, Female, Other:
		return true
	default:
		return false
	}
}
