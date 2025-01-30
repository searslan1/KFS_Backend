package auth

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/smtp"
	"time"
	"os"
)

type AuthService struct {
	Repo *AuthRepository
}

// Kullanıcı kaydı
func (s *AuthService) RegisterUser(email, password, userType string) error {
	// Şifre hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		UserType:     userType,
		CreatedAt:    time.Now(),
	}

	return s.Repo.CreateUser(user)
}

// Kullanıcı giriş doğrulama
func (s *AuthService) AuthenticateUser(email, password string) (*User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("kullanıcı bulunamadı")
	}

	// Şifre doğrulama
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("şifre yanlış")
	}

	return user, nil
}

func (s *AuthService) GetUserByID(userID int64) (*User, error) {
	return s.Repo.GetUserByID(userID)
}
func (s *AuthService) SendEmailOTP(email string) error {
	user, err := s.Repo.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// OTP oluştur (6 haneli)
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	expireTime := time.Now().Add(15 * time.Minute)

	// OTP bilgilerini kaydet
	user.OTP = otp
	user.OTPExpiresAt = expireTime
	if err := s.Repo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// Kullanıcının e-posta adresine OTP gönder
	if err := sendEmailWithSMTP(email, otp); err != nil {
		return fmt.Errorf("failed to send email OTP: %w", err)
	}

	return nil
}

func (s *AuthService) SendPhoneOTP(phoneNumber string) error {
	user, err := s.Repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// OTP oluştur (6 haneli)
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	expireTime := time.Now().Add(15 * time.Minute)

	// OTP bilgilerini kaydet
	user.OTP = otp
	user.OTPExpiresAt = expireTime
	if err := s.Repo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// OTP'yi SMS ile gönder (örnek: Twilio kullanılabilir)
	fmt.Printf("Telefon OTP: %s\n", otp) // Test ortamında OTP'yi konsola yaz
	return nil
}
func sendEmailWithSMTP(toEmail, otp string) error {
	// Gönderen (sender) e-posta adresi ve şifresi
	senderEmail := os.Getenv("EMAIL_SENDER")
	senderPassword := os.Getenv("EMAIL_PASSWORD")

	// SMTP bilgileri
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// E-posta içeriği
	subject := "Subject: Your OTP Code\n"
	body := fmt.Sprintf("Your OTP code is: %s", otp)
	message := []byte(subject + "\n" + body)

	// SMTP kimlik doğrulama
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// E-postayı alıcıya gönder
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{toEmail}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
