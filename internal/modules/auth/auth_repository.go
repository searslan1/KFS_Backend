package auth

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"KFS_Backend/internal/utils"
)

// AuthRepository, veritabanı işlemleri için kullanılan repository yapısıdır.
type AuthRepository struct {
	DB *gorm.DB // Veritabanı bağlantısı
}

// Kullanıcı oluşturma işlemi. User tablosuna yeni bir kullanıcı ekler.
func (r *AuthRepository) CreateUser(user *User) error {
	// GORM'un Create fonksiyonunu kullanarak User tablosuna veri ekler.
	return r.DB.Create(user).Error
}

// E-posta adresine göre kullanıcıyı veritabanından getirir.
func (r *AuthRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	// E-posta adresine göre kullanıcıyı sorgular.
	err := r.DB.Where("email = ?", email).First(&user).Error
	// Eğer kullanıcı bulunamazsa hata döner.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

// Kullanıcı ID'sine göre kullanıcıyı veritabanından getirir.
func (r *AuthRepository) GetUserByID(userID int64) (*User, error) {
	var user User
	// Kullanıcı ID'sine göre kullanıcıyı sorgular.
	err := r.DB.Where("user_id = ?", userID).First(&user).Error
	// Eğer kullanıcı bulunamazsa hata döner.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, err
}

// Başarısız giriş denemelerini günceller.
func (r *AuthRepository) UpdateFailedLoginAttempts(userID int64, attempts int) error {
	// Kullanıcıya ait failed_login_attempts alanını günceller.
	return r.DB.Model(&AuthUser{}).Where("user_id = ?", userID).Update("failed_login_attempts", attempts).Error
}

// Verification tablosuna yeni bir doğrulama kodu kaydeder.
func (r *AuthRepository) SaveVerificationCode(verification *utils.Verification) error {
	// Veritabanına doğrulama kodu ekler.
	err := r.DB.Create(verification).Error
	if err != nil {
		return errors.New("failed to save verification code")
	}
	return nil
}

// Kullanıcı ID'si ve doğrulama türüne göre Verification kaydını getirir.
func (r *AuthRepository) GetVerificationCode(userID int64, verificationType string) (*utils.Verification, error) {
	var verification utils.Verification
	// Kullanıcı ID'si ve türüne göre doğrulama kodunu sorgular.
	err := r.DB.Where("user_id = ? AND type = ?", userID, verificationType).First(&verification).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("verification code not found")
		}
		return nil, err
	}
	return &verification, nil
}

// AuthUser tablosunda bir kullanıcıyı günceller.
func (r *AuthRepository) UpdateAuthUser(authUser *AuthUser) error {
	// Verilen AuthUser nesnesini günceller.
	err := r.DB.Save(authUser).Error
	if err != nil {
		return fmt.Errorf("failed to update auth user: %w", err)
	}
	return nil
}

// Kullanıcı ID'sine göre AuthUser tablosundan bir kullanıcı getirir.
func (r *AuthRepository) GetAuthUserByID(userID int64) (*AuthUser, error) {
	var authUser AuthUser
	// Kullanıcı ID'sine göre AuthUser kaydını sorgular.
	err := r.DB.Where("user_id = ?", userID).First(&authUser).Error
	// Eğer kayıt bulunamazsa hata döner.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("auth user not found")
	}
	return &authUser, err
}

// Verification kaydını günceller.
func (r *AuthRepository) UpdateVerificationCode(verification *utils.Verification) error {
	// Verilen Verification nesnesini günceller.
	err := r.DB.Save(verification).Error
	if err != nil {
		return fmt.Errorf("failed to update verification code: %w", err)
	}
	return nil
}

// AuthUser tablosuna yeni bir kullanıcı kaydı ekler.
func (r *AuthRepository) CreateAuthUser(authUser *AuthUser) error {
	// AuthUser tablosuna yeni bir kayıt ekler.
	err := r.DB.Create(authUser).Error
	if err != nil {
		return fmt.Errorf("failed to create auth user: %w", err)
	}
	return nil
}
func (r *AuthRepository) GetAllUsers() ([]User, error) {
	var users []User
	// Tüm kullanıcıları çekmek için sorgu
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}
