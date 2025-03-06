package campaign

import (
	"errors"

	"gorm.io/gorm"
)

// Repository interface, campaign modülü ile ilgili tüm veritabanı işlemlerini tanımlar.
type Repository interface {
	CreateCampaignProfile(campaign *CampaignProfile) error
	GetCampaignProfileByID(id uint) (*CampaignProfile, error)
	UpdateCampaignProfile(campaign *CampaignProfile) error
	DeleteCampaignProfile(id uint) error
}

// repository yapısı, GORM DB nesnesini içerir.
type repository struct {
	db *gorm.DB
}

// NewRepository, campaign modülü için yeni bir repository örneği oluşturur.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateCampaignProfile, yeni bir kampanya profili oluşturur.
func (r *repository) CreateCampaignProfile(campaign *CampaignProfile) error {
	return r.db.Create(campaign).Error
}

// GetCampaignProfileByID, belirtilen ID'ye sahip kampanya profilini ilişkili alanları preload ederek getirir.
func (r *repository) GetCampaignProfileByID(id uint) (*CampaignProfile, error) {
	var campaign CampaignProfile
	if err := r.db.
		Preload("Prizes").
		Preload("Patents").
		Preload("Law").
		Preload("TeamMembers").
		Preload("ProductProduction").
		Preload("Risks").
		Preload("Funding").
		Preload("Establishment").
		Preload("FinancialTables").
		Preload("VisualVideos").
		Preload("OtherDocuments").
		First(&campaign, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Kayıt bulunamazsa nil döndür.
		}
		return nil, err
	}
	return &campaign, nil
}

// UpdateCampaignProfile, mevcut kampanya profilini günceller.
func (r *repository) UpdateCampaignProfile(campaign *CampaignProfile) error {
	return r.db.Save(campaign).Error
}

// DeleteCampaignProfile, belirtilen ID'ye sahip kampanya profilini siler.
func (r *repository) DeleteCampaignProfile(id uint) error {
	return r.db.Delete(&CampaignProfile{}, id).Error
}
