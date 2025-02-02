package utils

import (
	"crypto/aes"           // AES (Advanced Encryption Standard) şifreleme algoritması için kullanılan kütüphane.
	"crypto/cipher"        // Şifreleme ve deşifreleme işlemleri için kullanılan kütüphane.
	"encoding/base64"      // Base64 formatında kodlama ve çözme işlemleri için kullanılan kütüphane.
)

// AES algoritması ile verilen veriyi şifreler.
func EncryptAES(data string, key string) (string, error) {
	// Şifreleme anahtarını (key) kullanarak bir AES şifreleme bloğu oluşturur.
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		// Şifreleme bloğu oluşturulamazsa bir hata döner.
		return "", err
	}

	// Şifrelenecek veriyi bayt dizisine dönüştürür.
	plaintext := []byte(data)

	// AES'in CFB (Cipher Feedback Mode) modunda bir şifreleme akışı oluşturur.
	// Burada `block.BlockSize()` şifreleme anahtarının boyutuna göre ivme (initialization vector) oluşturmak için kullanılır.
	cfb := cipher.NewCFBEncrypter(block, []byte(key)[:block.BlockSize()])
	
	// Şifrelenmiş veriyi depolamak için bir bayt dizisi oluşturur.
	ciphertext := make([]byte, len(plaintext))

	// Şifreleme işlemi yapılır. Verinin şifrelenmiş hali `ciphertext` değişkenine yazılır.
	cfb.XORKeyStream(ciphertext, plaintext)

	// Şifrelenmiş veriyi Base64 formatına çevirerek döner.
	// Base64, ikili verileri metin formatında saklamak ve taşımak için kullanılır.
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES algoritması ile şifrelenmiş veriyi çözerek orijinal haline getirir.
func DecryptAES(encryptedData string, key string) (string, error) {
	// Şifreleme anahtarını (key) kullanarak bir AES şifreleme bloğu oluşturur.
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		// Şifreleme bloğu oluşturulamazsa bir hata döner.
		return "", err
	}
	// Base64 formatındaki şifrelenmiş veriyi çözer ve bayt dizisine dönüştürür.
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		// Base64 çözme işlemi başarısız olursa bir hata döner.
		return "", err
	}
	// AES'in CFB modunda bir deşifreleme akışı oluşturur.
	cfb := cipher.NewCFBDecrypter(block, []byte(key)[:block.BlockSize()])
	
	// Çözülen veriyi depolamak için bir bayt dizisi oluşturur.
	plaintext := make([]byte, len(ciphertext))

	// Deşifreleme işlemi yapılır. Orijinal veri `plaintext` değişkenine yazılır.
	cfb.XORKeyStream(plaintext, ciphertext)

	// Çözülen veriyi string formatında döner.
	return string(plaintext), nil
}
