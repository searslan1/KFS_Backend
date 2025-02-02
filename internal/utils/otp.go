package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// OTP oluşturma fonksiyonu
func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // 6 haneli OTP
}

// OTP'nin geçerlilik süresini kontrol et
func IsOTPValid(expiry time.Time) bool {
	return time.Now().Before(expiry)
}
