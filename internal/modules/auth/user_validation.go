package user

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("geçersiz e-posta adresi")
	}
	return nil
}
