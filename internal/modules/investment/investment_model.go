package investment

type Investment struct {
	ID         int `json:"id"`
	UserID     int `json:"user_id"`
	CampaignID int `json:"campaign_id"`
	Amount     int `json:"amount"`
}
