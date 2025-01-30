package auth

import (
	"errors"
	"gorm.io/gorm"
	//"time"
	"fmt"
)

type AuthRepository struct {
	DB *gorm.DB
}

// Kullanıcı oluşturma
func (r *AuthRepository) CreateUser(user *User) error {
	// GORM'un Create fonksiyonunu kullanıyoruz
	return r.DB.Create(user).Error
}

// Email ile kullanıcıyı getir
func (r *AuthRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	// GORM'un Where ve First fonksiyonlarını kullanarak kullanıcıyı alıyoruz
	err := r.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *AuthRepository) GetUserByID(userID int64) (*User, error) {
	var user User
	// GORM'un Where ve First fonksiyonlarını kullanarak kullanıcıyı alıyoruz
	err := r.DB.Where("user_id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

// Başarısız giriş denemelerini güncelle
func (r *AuthRepository) UpdateFailedLoginAttempts(userID int64, attempts int) error {
	// GORM'un Update fonksiyonlarını kullanıyoruz
	return r.DB.Model(&AuthUser{}).Where("user_id = ?", userID).Update("failed_login_attempts", attempts).Error
}

func (r *AuthRepository) UpdateUser(user *User) error {
	return r.DB.Save(user).Error
}

func (r *AuthRepository) GetUserByPhoneNumber(phoneNumber string) (*User, error) {
	var user User
	err := r.DB.Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}
