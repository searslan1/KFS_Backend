package campaign

import (
	"errors"
)

// Service arayüzü, campaign modülü için temel işlemleri tanımlar.
type Service interface {
	CreateCampaign(campaign *CampaignProfile) error
	GetCampaignByID(id uint) (*CampaignProfile, error)
	UpdateCampaign(campaign *CampaignProfile) error
	DeleteCampaign(id uint) error
}

// service yapısı, repository arayüzüne bağımlıdır.
type service struct {
	repo Repository
}

// NewService, campaign modülü için yeni bir service örneği oluşturur.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateCampaign, yeni bir kampanya profili oluşturur.
func (s *service) CreateCampaign(campaign *CampaignProfile) error {
	return s.repo.CreateCampaignProfile(campaign)
}

// GetCampaignByID, verilen ID'ye ait kampanya profilini getirir.
func (s *service) GetCampaignByID(id uint) (*CampaignProfile, error) {
	campaign, err := s.repo.GetCampaignProfileByID(id)
	if err != nil {
		return nil, err
	}
	if campaign == nil {
		return nil, errors.New("campaign not found")
	}
	return campaign, nil
}

// UpdateCampaign, varolan kampanya profilini günceller.
func (s *service) UpdateCampaign(campaign *CampaignProfile) error {
	return s.repo.UpdateCampaignProfile(campaign)
}

// DeleteCampaign, belirtilen ID'ye sahip kampanya profilini siler.
func (s *service) DeleteCampaign(id uint) error {
	return s.repo.DeleteCampaignProfile(id)
}
