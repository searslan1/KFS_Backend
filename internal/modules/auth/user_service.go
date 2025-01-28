package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserService struct {
	Repo *UserRepository
}

// Kullanıcı kaydı
func (s *UserService) RegisterUser(email, password, userType string) error {
	// Şifre hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Email:     email,
		Password:  string(hashedPassword),
		UserType:  userType,
		CreatedAt: time.Now(),
	}

	return s.Repo.CreateUser(user)
}

// Kullanıcı giriş doğrulama
func (s *UserService) AuthenticateUser(email, password string) (*User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	// Şifre doğrulama
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("şifre yanlış")
	}

	return user, nil
}
