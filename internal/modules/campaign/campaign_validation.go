package campaign

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateCampaign, kampanya verisinin doğrulanması için kullanılır.
func ValidateCampaign(campaign *CampaignProfile) error {
	return validate.Struct(campaign)
}
