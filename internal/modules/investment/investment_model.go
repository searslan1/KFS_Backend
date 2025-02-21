package investment

import "time"

// Investment represents a campaign investment record
type Investment struct {
	ID             uint      `json:"id" gorm:"primaryKey;column:investment_id"`
	CampaignID     uint      `json:"campaign_id" gorm:"not null"`
	InvestorID     uint      `json:"investor_id" gorm:"not null"`
	Amount         float64   `json:"amount" gorm:"type:decimal(15,2);not null"`
	InvestmentDate time.Time `json:"investment_date" gorm:"not null"`
	Status         string    `json:"status" gorm:"type:varchar(20)"`
}

// // InvestorProfile represents an investor's profile
// type InvestorProfile struct {
// 	ID                  uint      `json:"id" gorm:"primaryKey;column:investor_id"`
// 	UserID              uint      `json:"user_id" gorm:"not null"`
// 	InvestmentFocus     string    `json:"investment_focus" gorm:"type:text"`
// 	MinInvestmentAmount float64   `json:"min_investment_amount" gorm:"type:decimal(15,2)"`
// 	MaxInvestmentAmount float64   `json:"max_investment_amount" gorm:"type:decimal(15,2)"`
// 	PastInvestments     string    `json:"past_investments" gorm:"type:text"`
// 	RiskProfile         string    `json:"risk_profile" gorm:"type:varchar(50)"`
// 	Status              string    `json:"status" gorm:"type:varchar(20)"`
// 	CreatedAt           time.Time `json:"created_at" gorm:"not null"`
// }

// TableName overrides the table name for Investment
func (Investment) TableName() string {
	return "campaign_investments"
}

// // TableName overrides the table name for InvestorProfile
// func (InvestorProfile) TableName() string {
// 	return "investor_profiles"
// }
