package profile

import (
	"errors"
	"regexp"
)

// Kullanıcı profili doğrulama
func ValidateUserProfile(profile *UserProfile) error {
	if profile.Bio == "" {
		return errors.New("bio boş olamaz")
	}

	if profile.WebsiteURL != "" {
		matched, _ := regexp.MatchString(`^https?://`, profile.WebsiteURL)
		if !matched {
			return errors.New("geçersiz website URL")
		}
	}

	return nil
}

// Adres doğrulama
func ValidateAddress(address *Address) error {
	if address.Country == "" {
		return errors.New("ülke boş olamaz")
	}
	if address.City == "" {
		return errors.New("şehir boş olamaz")
	}
	return nil
}

// Rol doğrulama
func ValidateRoleProfile(role *RoleProfile) error {
	if role.CompanyName == "" {
		return errors.New("şirket adı boş olamaz")
	}
	if role.TaxID == "" {
		return errors.New("vergi kimlik numarası boş olamaz")
	}
	return nil
}
