package profile

import (
	"errors"
)

type ProfileService struct {
	repo *ProfileRepository
}

func NewProfileService(repo *ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// Kullanıcının profil bilgilerini getir
func (s *ProfileService) GetUserProfile(userID int64) (*UserProfile, error) {
	return s.repo.GetProfileByUserID(userID)
}

// Kullanıcının adres bilgilerini getir
func (s *ProfileService) GetUserAddress(userID int64) (*Address, error) {
	return s.repo.GetAddressByUserID(userID)
}

// Kullanıcının medya dosyalarını getir
func (s *ProfileService) GetUserMediaFiles(userID int64) ([]MediaFile, error) {
	return s.repo.GetMediaFilesByUserID(userID)
}

// Kullanıcının rol bilgilerini getir
func (s *ProfileService) GetUserRole(userID int64) (*RoleProfile, error) {
	return s.repo.GetRoleByUserID(userID)
}

// Kullanıcının profilini güncelle
func (s *ProfileService) UpdateUserProfile(userID int64, updatedProfile *UserProfile) error {
	if userID == 0 {
		return errors.New("geçersiz kullanıcı ID")
	}
	return s.repo.InsertOrUpdateProfile(updatedProfile)
}

// Kullanıcının adresini güncelle
func (s *ProfileService) UpdateUserAddress(userID int64, updatedAddress *Address) error {
	if userID == 0 {
		return errors.New("geçersiz kullanıcı ID")
	}
	return s.repo.InsertOrUpdateAddress(updatedAddress)
}

// Kullanıcının rol bilgisini güncelle
func (s *ProfileService) UpdateUserRole(userID int64, updatedRole *RoleProfile) error {
	if userID == 0 {
		return errors.New("geçersiz kullanıcı ID")
	}
	return s.repo.InsertOrUpdateRole(updatedRole)
}
