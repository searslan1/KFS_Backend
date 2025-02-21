package profile

import (
	"gorm.io/gorm"
	"log"
	//"KFS_Backend/internal/modules/profile"
)

// ProfileRepository yapısı, veritabanı işlemleri için kullanılan repository yapısıdır.
type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	if db == nil {
		panic("DB instance is nil")
	}
	return &ProfileRepository{db}
}

// Kullanıcının profil bilgilerini getirir.
func (r *ProfileRepository) GetProfileByUserID(userID int64) (*UserProfile, error) {
	var userProfile UserProfile
	if err := r.db.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

// Kullanıcının adres bilgilerini getirir.
func (r *ProfileRepository) GetAddressByUserID(userID int64) (*Address, error) {
	var address Address
	if err := r.db.Where("user_id = ?", userID).First(&address).Error; err != nil {
		return nil, err
	}
	return &address, nil
}

// Kullanıcının medya dosyalarını getirir.
func (r *ProfileRepository) GetMediaFilesByUserID(userID int64) ([]MediaFile, error) {
	var mediaFiles []MediaFile
	if err := r.db.Where("user_id = ?", userID).Find(&mediaFiles).Error; err != nil {
		return nil, err
	}
	return mediaFiles, nil
}

// Kullanıcının rol bilgilerini getirir.
func (r *ProfileRepository) GetRoleByUserID(userID int64) (*RoleProfile, error) {
	var role RoleProfile
	if err := r.db.Where("user_id = ?", userID).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// InsertOrUpdateProfile: Kayıt varsa günceller, yoksa ekler
func (repo *ProfileRepository) InsertOrUpdateProfile(profile *UserProfile) error {
	var existingProfile UserProfile
	err := repo.db.Where("user_id = ?", profile.UserID).First(&existingProfile).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Kayıt bulunamadı, yeni bir kayıt ekle
			log.Println("Profile not found. Creating a new one.")
			err = repo.db.Create(profile).Error
			if err != nil {
				return err
			}
			return nil
		}
		// Diğer hatalar için geri dön
		return err
	}

	// Kayıt bulundu, güncelle
	log.Println("Profile found. Updating existing profile.")
	existingProfile.ProfilePictureURL = profile.ProfilePictureURL
	existingProfile.Bio = profile.Bio
	existingProfile.WebsiteURL = profile.WebsiteURL
	existingProfile.SocialMedia = profile.SocialMedia
	err = repo.db.Save(&existingProfile).Error
	if err != nil {
		return err
	}

	return nil
}
func (r *ProfileRepository) InsertOrUpdateAddress(address *Address) error {
	var existingAddress Address
	err := r.db.Where("user_id = ?", address.UserID).First(&existingAddress).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Kayıt yoksa ekle
			log.Println("Address not found. Creating a new one.")
			return r.db.Create(address).Error
		}
		return err
	}

	// Kayıt varsa güncelle
	log.Println("Address found. Updating existing address.")
	existingAddress.Country = address.Country
	existingAddress.City = address.City
	existingAddress.District = address.District
	existingAddress.Street = address.Street
	existingAddress.ZipCode = address.ZipCode
	return r.db.Save(&existingAddress).Error
}


// Kullanıcının rol bilgisini ekler veya günceller.
func (r *ProfileRepository) InsertOrUpdateRole(role *RoleProfile) error {
	var existingRole RoleProfile
	err := r.db.Where("user_id = ?", role.UserID).First(&existingRole).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Kayıt bulunamadı, yeni bir kayıt ekle
			log.Println("Role not found. Creating a new one.")
			err = r.db.Create(role).Error
			if err != nil {
				return err
			}
			return nil
		}
		// Diğer hatalar için geri dön
		return err
	}

	// Kayıt bulundu, güncelle
	log.Println("Role found. Updating existing role.")
	existingRole.CompanyName = role.CompanyName
	existingRole.TaxID = role.TaxID
	existingRole.RegistrationNumber = role.RegistrationNumber
	existingRole.TaxOffice = role.TaxOffice
	existingRole.PhoneNumber = role.PhoneNumber
	err = r.db.Save(&existingRole).Error
	if err != nil {
		return err
	}

	return nil
}
