package validator

import (
	"regexp"
)

const (
	usernamePattern = `^[a-zA-Z][a-zA-Z0-9_-]$`
	emailPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]$`
	passwordPattern = `^[a-zA-Z][a-zA-Z0-9_-]$`
	//passwordPattern = `^(?=.*[a-zA-Z])(?=.*\d)(?=.*[@!#%^&*])[a-zA-Z\d@!#%^&*]$`
)

func IsLengthValid(str string, min, max int) bool {
	length := len(str)
	return length >= min && length <= max
}

func ValidateByPattern(pattern, s string) bool {
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}

func ValidateUsername(username string) bool {
	return ValidateByPattern(usernamePattern, username)
}

func ValidateEmail(email string) bool {
	return ValidateByPattern(emailPattern, email)
}

func ValidatePassword(password string) bool {
	return ValidateByPattern(passwordPattern, password)
}
