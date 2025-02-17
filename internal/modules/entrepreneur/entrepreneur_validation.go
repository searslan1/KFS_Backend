package validation

import (
	"errors"
	"entrepreneur/model"
)

// ValidateEntrepreneur girişimci verilerini doğrular
func ValidateEntrepreneur(e *model.Entrepreneur) error {
	if e.UserID == 0 {
		return errors.New("UserID boş olamaz ve zorunludur")
	}

	if e.StartupName == "" {
		return errors.New("Startup (Girişim) adı boş olamaz ve zorunludur")
	}

	if len(e.StartupName) < 3 {
		return errors.New("Startup adı en az 3 karakter olmalıdır")
	}

	if e.Industry == "" {
		return errors.New("Sektör bilgisi boş olamaz ve zorunludur")
	}

	if e.FundingNeeded <= 0 {
		return errors.New("Yatırım miktarı sıfırdan büyük olmalıdır")
	}

	if e.BusinessModel == "" {
		return errors.New("İş modeli boş olamaz ve zorunludur")
	}

	if e.PitchDeckURL == "" {
		return errors.New("Pitch Deck (Sunum) URL'si boş olamaz ve zorunludur")
	}

	if !e.IsEDevletApproved {
		return errors.New("Girişimci profili oluşturmak için e-Devlet onayı gereklidir")
	}

	return nil
}

// ValidateAdminApproval girişimci profili admin tarafından onaylanırken doğrulama yapar
func ValidateAdminApproval(e *model.Entrepreneur) error {
	if e.Status == "approved" {
		return errors.New("Bu girişimci zaten onaylanmış durumda")
	}

	if !e.IsEDevletApproved {
		return errors.New("E-Devlet onayı yapılmadan girişimci profili onaylanamaz")
	}

	return nil
}
