package utils

import (
	"fmt"
	"regexp"
)

type InvalidEmailRegexError struct {
	Email string
}

func (e *InvalidEmailRegexError) Error() string {
	return fmt.Sprintf("Invalid email format: %s", e.Email)
}

func EmailVerification(email string) (string, error) {
	regex := "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	re := regexp.MustCompile(regex)

	if re.MatchString(email) {
		return email, nil
	}

	return "", &InvalidEmailRegexError{Email: email}
}
