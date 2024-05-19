package domain_models

type UserGender string

const (
	Male   UserGender = "Male"
	Female UserGender = "Female"
	Other  UserGender = "Other"
)

// IsValidGender function to check if the given value is a valid UserGender.
func IsValidGender(gender UserGender) bool {
	switch gender {
	case Male, Female, Other:
		return true
	default:
		return false
	}
}

// GetGender returns the string representation of a UserGender enum value.
func GetGender(gender UserGender) string {
	switch gender {
	case Male:
		return "Male"
	case Female:
		return "Female"
	case Other:
		return "Other"
	default:
		return ""
	}
}

// GetGenderType returns a UserGender enum value.
func GetGenderType(gender string) UserGender {
	switch gender {
	case "Male":
		return Male
	case "Female":
		return Female
	case "Other":
		return Other
	default:
		return Other
	}
}

// GetGenderTypeInt returns a UserGender enum value.
func GetGenderTypeInt(gender int) UserGender {
	switch gender {
	case 0:
		return Other
	case 1:
		return Female
	case 2:
		return Male
	default:
		return Other
	}
}
