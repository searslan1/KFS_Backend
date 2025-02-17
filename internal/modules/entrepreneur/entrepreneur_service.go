package service

import (
	"entrepreneur/model"
	"errors"
)

// EntrepreneurService girişimci işlemlerini yönetir
type EntrepreneurService struct {
	repo *model.EntrepreneurRepository
}

// NewEntrepreneurService yeni bir hizmet oluşturur
func NewEntrepreneurService(repo *model.EntrepreneurRepository) *EntrepreneurService {
	return &EntrepreneurService{repo: repo}
}

// IsEDevletVerified e-Devlet doğrulamasını kontrol eder
func (s *EntrepreneurService) IsEDevletVerified(userID uint) bool {
	// SPK Gerekliliği: Burada e-Devlet entegrasyonu ile doğrulama yapılmalıdır
	// Gerçek entegrasyon sağlandığında, bu fonksiyonun döndüreceği değer gerçek veriye dayanmalıdır
	return true // Şimdilik varsayılan olarak true dönüyor
}

// GetByUserID belirli bir kullanıcının girişimci profilini getirir
func (s *EntrepreneurService) GetByUserID(userID uint) (*model.Entrepreneur, error) {
	return s.repo.GetByUserID(userID)
}

// Create yeni girişimci profili oluşturur
func (s *EntrepreneurService) Create(entrepreneur *model.Entrepreneur) error {
	// Kullanıcının zaten bir girişimci profili olup olmadığını kontrol et
	existing, _ := s.repo.GetByUserID(entrepreneur.UserID)
	if existing != nil {
		return errors.New("Kullanıcı zaten bir girişimci profiline sahip")
	}

	// Girişimci profilini oluştur
	return s.repo.Create(entrepreneur)
}

// Update girişimci profilini günceller
func (s *EntrepreneurService) Update(entrepreneur *model.Entrepreneur) error {
	return s.repo.Update(entrepreneur)
}
