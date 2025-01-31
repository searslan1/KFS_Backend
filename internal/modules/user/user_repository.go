package user

import (
	"fmt"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateUser(user *User) error {
	tx := r.db.Begin()

	// 1. Kullanıcıyı oluştur ve ID'nin atanmasını bekle
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Kullanıcı ID'sinin atandığını doğrula
	if user.UserID == 0 {
		tx.Rollback()
		return fmt.Errorf("user ID atanamadı, transaction iptal edildi")
	}

	// 2. Kullanıcı doğrulama kaydını oluştur
	authUser := AuthUser{
		UserID: user.UserID,
	}
	if err := tx.Create(&authUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. Commit işlemi
	return tx.Commit().Error
}

func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) SaveUserSession(session *UserSession) error {
	return r.db.Create(session).Error
}

func (r *UserRepository) GetSessionByRefreshToken(refreshToken string) (*UserSession, error) {
	var session UserSession
	err := r.db.Where("refresh_token = ?", refreshToken).First(&session).Error
	return &session, err
}
