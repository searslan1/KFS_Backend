package user

import (
	"time"

	"KFS_Backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) RegisterUser(email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Email:        email,
		PasswordHash: string(hashedPassword),
	}
	err = s.repo.CreateUser(&user)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(email, password, ip, deviceID string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", err
	}

	accessToken, err := utils.GenerateAccessToken(user.UserID, user.Email, user.Role, ip, deviceID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateSecureRefreshToken()
	if err != nil {
		return "", "", err
	}

	session := UserSession{
		UserID:             user.UserID,
		IPAddress:          ip,
		DeviceInfo:         deviceID,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: time.Now().Add(7 * 24 * time.Hour),
	}
	s.repo.SaveUserSession(&session)

	return accessToken, refreshToken, nil
}
