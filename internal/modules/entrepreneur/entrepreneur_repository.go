package repository

import (
	"entrepreneur/model"
	"errors"
	"gorm.io/gorm"
)

// EntrepreneurRepository girişimci profili ile ilgili veritabanı işlemlerini yönetir
type EntrepreneurRepository struct {
	db *gorm.DB
}

// NewEntrepreneurRepository yeni bir EntrepreneurRepository örneği oluşturur
func NewEntrepreneurRepository(db *gorm.DB) *EntrepreneurRepository {
	return &EntrepreneurRepository{db: db}
}

// Create yeni bir girişimci profili ekler
func (r *EntrepreneurRepository) Create(entrepreneur *model.Entrepreneur) error {
	// Kullanıcı daha önce girişimci profili oluşturmuş mu kontrol et
	var existingProfile model.Entrepreneur
	err := r.db.Where("user_id = ?", entrepreneur.UserID).First(&existingProfile).Error
	if err == nil {
		return errors.New("kullanıcı zaten bir girişimci profili oluşturmuş")
	}

	// E-Devlet onayı kontrolü
	if !entrepreneur.IsEDevletApproved {
		return errors.New("e-Devlet onayı gereklidir")
	}

	// Yeni girişimci profilini oluştur
	return r.db.Create(entrepreneur).Error
}

// GetByID girişimci profilini ID'ye göre getirir
func (r *EntrepreneurRepository) GetByID(id uint) (*model.Entrepreneur, error) {
	var entrepreneur model.Entrepreneur
	err := r.db.First(&entrepreneur, id).Error
	return &entrepreneur, err
}

// GetByUserID girişimci profilini user_id'ye göre getirir
func (r *EntrepreneurRepository) GetByUserID(userID uint) (*model.Entrepreneur, error) {
	var entrepreneur model.Entrepreneur
	err := r.db.Where("user_id = ?", userID).First(&entrepreneur).Error
	return &entrepreneur, err
}

// Update girişimci profilini günceller
func (r *EntrepreneurRepository) Update(entrepreneur *model.Entrepreneur) error {
	// Admin onayını kontrol et ve statü güncellemesi yap
	if entrepreneur.IsAdminApproved {
		entrepreneur.Status = "approved"
	}
	return r.db.Save(entrepreneur).Error
}

// Delete belirtilen ID'ye sahip girişimci profilini siler
func (r *EntrepreneurRepository) Delete(id uint) error {
	return r.db.Delete(&model.Entrepreneur{}, id).Error
}

// CheckForExistingProfile kullanıcıya ait girişimci profili olup olmadığını kontrol eder
func (r *EntrepreneurRepository) CheckForExistingProfile(userID uint) bool {
	var entrepreneur model.Entrepreneur
	err := r.db.Where("user_id = ?", userID).First(&entrepreneur).Error
	return err == nil
}
