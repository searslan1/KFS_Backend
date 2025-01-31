package campaign

import (
	"time"
)

// CampaignProfile tablosu modeli
type Campaign struct {
	ID                    uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                uint           `gorm:"not null;index" json:"user_id"` // users tablosuna referans
	CampaignLogo          string         `gorm:"type:text" json:"campaign_logo"`
	EntrepreneurName      string         `gorm:"type:varchar(255)" json:"entrepreneur_name"`
	CampaignName          string         `gorm:"type:varchar(255);not null" json:"campaign_name"`
	CampaignDescription   string         `gorm:"type:text" json:"campaign_description"`
	AboutProject          string         `gorm:"type:text" json:"about_project"`
	CampaignSummary       string         `gorm:"type:text" json:"campaign_summary"`
	GoalCoverageSubject   string         `gorm:"type:text" json:"goal_coverage_subject"`
	EntrepreneurStageID   uint           `gorm:"index" json:"entrepreneur_stage_id"`
	Location              string         `gorm:"type:varchar(255)" json:"location"`
	Category              string         `gorm:"type:varchar(100)" json:"category"`
	BusinessModelsID      uint           `gorm:"index" json:"business_models_id"`
	Sector                string         `gorm:"type:varchar(100)" json:"sector"`
	IsPastCampaign        bool           `gorm:"default:false" json:"is_past_campaign"`
	CampaignStatus        string         `gorm:"type:varchar(50);default:'pending'" json:"campaign_status"` // pending, approved, rejected, active, completed
	CreatedAt             time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	EntrepreneursMails    []CampaignEntrepreneur `gorm:"foreignKey:CampaignID" json:"entrepreneurs_mails"`
}

// CampaignEntrepreneur tablosu (Girişimcilerin e-posta adresleri için ayrı tablo)
type CampaignEntrepreneur struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID uint   `gorm:"not null;index" json:"campaign_id"` // CampaignProfile ile ilişkilendirme
	Email      string `gorm:"type:varchar(255);not null" json:"email"`
}

// GORM'un CampaignEntrepreneurs tablosunu "campaign_entrepreneurs" olarak kullanmasını sağlıyoruz
func (CampaignEntrepreneur) TableName() string {
	return "campaign_entrepreneurs"
}
