package auth

import (
	"KFS_Backend/internal/utils"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type AuthService struct {
	Repo *AuthRepository // Veritabanı işlemleri için kullanılan repository
}

// Yeni bir kullanıcı kaydı oluşturur.
func (s *AuthService) RegisterUser(email, password, userType string) error {
	// Kullanıcının şifresini hashler.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Kullanıcının e-posta adresini şifreler.
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	encryptedEmail, err := utils.EncryptAES(email, encryptionKey)
	if len(encryptionKey) != 16 && len(encryptionKey) != 24 && len(encryptionKey) != 32 {
		return fmt.Errorf("invalid encryption key size: must be 16, 24, or 32 bytes")
	}

	// Yeni kullanıcı nesnesi oluşturur.
	user := &User{
		Email:        encryptedEmail,
		PasswordHash: string(hashedPassword),
		UserType:     userType,
		CreatedAt:    time.Now(),
	}

	// Kullanıcıyı veritabanına kaydeder.
	if err := s.Repo.CreateUser(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// AuthUser kaydını oluşturur.
	authUser := &AuthUser{
		UserID:        int64(user.UserID),
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:     time.Now(),
	}

	// AuthUser tablosuna kaydeder.
	if err := s.Repo.CreateAuthUser(authUser); err != nil {
		return fmt.Errorf("failed to create auth user: %w", err)
	}

	return nil
}

func (s *AuthService) AuthenticateUser(email, password, ip, userAgent string) (string, string, error) {
	// Kullanıcıyı getir
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	encryptedEmail, err := utils.EncryptAES(email, encryptionKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to encrypt email: %w", err)
	}

	user, err := s.Repo.GetUserByEmail(encryptedEmail)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	// Şifre doğrulama işlemi
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid password")
	}

	// ✅ **Access Token oluştur**
	accessToken, err := utils.GenerateAccessToken(uint(user.UserID), user.UserType)
	if err != nil {
		return "", "", err
	}

	// ✅ **Refresh Token oluştur**
	refreshToken, err := utils.GenerateSecureRefreshToken()
	if err != nil {
		return "", "", err
	}

	// ✅ **Refresh Token kaydet**
	session := UserSession{
		UserID:             uint(user.UserID),
		IPAddress:          ip,
		DeviceInfo:         userAgent,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: time.Now().Add(7 * 24 * time.Hour), // 7 gün geçerli
	}
	err = s.Repo.SaveUserSession(&session)
	if err != nil {
		log.Println("❌ Refresh token veritabanına kaydedilemedi:", err)
	}

	return accessToken, refreshToken, nil
}

// Kullanıcıyı ID'ye göre getirir.
func (s *AuthService) GetUserByID(userID int64) (*User, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Kullanıcının e-posta adresini deşifre eder.
	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	decryptedEmail, err := utils.DecryptAES(user.Email, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt email: %w", err)
	}
	user.Email = decryptedEmail

	return user, nil
}

// Kullanıcıya e-posta doğrulama kodu gönderir.
func (s *AuthService) SendEmailVerification(userID int64, email string) error {
	// Rastgele 6 haneli bir OTP oluşturur.
	otp := utils.GenerateOTP()
	expiry := time.Now().Add(15 * time.Minute)

	// Doğrulama kodu kaydı oluşturur.
	verification := &utils.Verification{
		UserID:     userID,
		Code:       otp,
		CodeExpiry: expiry,
		Type:       "email",
	}

	// Doğrulama kodunu veritabanına kaydeder.
	if err := s.Repo.SaveVerificationCode(verification); err != nil {
		return fmt.Errorf("failed to save verification code: %w", err)
	}

	// Kullanıcıya e-posta gönderir.
	subject := "Your Email Verification Code"
	body := fmt.Sprintf("Your email verification code is: %s", otp)
	return utils.SendEmailWithSendGrid(email, subject, body)
}

// Kullanıcıdan gelen OTP'yi doğrular.
func (s *AuthService) VerifyEmailOTP(userID int64, otp string) error {
	// Veritabanından doğrulama kaydını alır.
	verification, err := s.Repo.GetVerificationCode(userID, "email")
	if err != nil {
		return fmt.Errorf("verification code not found: %w", err)
	}

	// Kullanıcının gönderdiği OTP ile eşleşip eşleşmediğini kontrol eder.
	if verification.Code != otp {
		return fmt.Errorf("invalid OTP")
	}

	// OTP'nin geçerlilik süresini kontrol eder.
	if time.Now().After(verification.CodeExpiry) {
		return fmt.Errorf("OTP expired")
	}

	// AuthUser kaydını doğrulandı olarak işaretler.
	authUser, err := s.Repo.GetAuthUserByID(userID)
	if err != nil {
		return fmt.Errorf("auth user not found: %w", err)
	}

	authUser.EmailVerified = true

	// AuthUser kaydını günceller.
	if err := s.Repo.UpdateAuthUser(authUser); err != nil {
		return fmt.Errorf("failed to update auth user verification status: %w", err)
	}

	// Doğrulama kaydını doğrulandı olarak günceller.
	verification.IsVerified = true
	if err := s.Repo.UpdateVerificationCode(verification); err != nil {
		return fmt.Errorf("failed to update verification status: %w", err)
	}

	return nil
}
func (s *AuthService) GetAllUsers() ([]User, error) {
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}
func (s *AuthService) RefreshAccessToken(refreshToken, ip, userAgent string) (string, string, error) {
    // ✅ 1️⃣ Önce Blacklist’te olup olmadığını kontrol et
    blacklisted, err := s.Repo.IsTokenBlacklisted(refreshToken)
    if err != nil {
        return "", "", fmt.Errorf("error checking blacklist: %w", err)
    }
    if blacklisted {
        return "", "", errors.New("refresh token is blacklisted")
    }

    // ✅ 2️⃣ Token’a bağlı oturumu getir
    session, err := s.Repo.GetSessionByRefreshToken(refreshToken)
    if err != nil || session == nil || time.Now().After(session.RefreshTokenExpiry) {
        return "", "", errors.New("invalid or expired refresh token")
    }

    // ✅ 3️⃣ IP & User-Agent doğrulama ekleyelim (User-Agent boşsa kontrol etme)
    if session.IPAddress != ip {
        return "", "", errors.New("token validation failed: IP mismatch")
    }
    if session.DeviceInfo != "" && session.DeviceInfo != userAgent {
        return "", "", errors.New("token validation failed: device mismatch")
    }

    // ✅ 4️⃣ Eski Refresh Token’ı Blacklist’e ekleyelim!
    err = s.Repo.AddToBlacklist(refreshToken)
    if err != nil {
        return "", "", fmt.Errorf("failed to blacklist token: %w", err)
    }

    // ✅ 5️⃣ Yeni Refresh Token oluştur
    newRefreshToken, err := utils.GenerateSecureRefreshToken()
    if err != nil {
        return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
    }

    // ✅ 6️⃣ Yeni Refresh Token’ı yeni bir session olarak kaydet (tekrar eden session_id hatasını önlemek için)
    newSession := &UserSession{
        UserID:             session.UserID,
        IPAddress:          ip,
        DeviceInfo:         userAgent,
        RefreshToken:       newRefreshToken,
        RefreshTokenExpiry: time.Now().Add(7 * 24 * time.Hour), // 7 gün geçerli
    }

    err = s.Repo.SaveUserSession(newSession)
    if err != nil {
        return "", "", fmt.Errorf("failed to save new session: %w", err)
    }

    return "newAccessToken", newRefreshToken, nil
}


func (s *AuthService) LogoutUser(refreshToken string) error {
	// Refresh Token'ın veritabanında olup olmadığını kontrol et
	session, err := s.Repo.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// Refresh Token'ı Blacklist tablosuna ekle
	blacklistedToken := BlacklistedToken{
		Token:  refreshToken,
		Expiry: session.RefreshTokenExpiry, // Token süresi dolana kadar saklanacak
	}

	err = s.Repo.AddToBlacklist(blacklistedToken.Token)
	if err != nil {
		return err
	}

	// Kullanıcının mevcut oturumunu veritabanından sil
	return s.Repo.DeleteUserSession(session.SessionID)
}
