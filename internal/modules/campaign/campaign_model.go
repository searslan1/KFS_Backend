package campaign

import (
	"time"
)

// CampaignProfile modeli
type CampaignProfile struct {
	ID                  uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID              uint             `gorm:"not null" json:"user_id"`
	CampaignLogo        string           `json:"campaign_logo"`
	EntrepreneurName    string           `gorm:"size:255" json:"entrepreneur_name"`
	CampaignName        string           `gorm:"size:255;not null" json:"campaign_name"`
	CampaignDescription string           `json:"campaign_description"`
	AboutProject        string           `json:"about_project"`
	CampaignSummary     string           `json:"campaign_summary"`
	GoalCoverageSubject string           `json:"goal_coverage_subject"`
	EntrepreneurStageID uint             `json:"entrepreneur_stage_id"`
	Location            string           `gorm:"size:255" json:"location"`
	Category            string           `gorm:"size:100" json:"category"`
	BusinessModelsID    uint             `json:"business_models_id"`
	Sector              string           `gorm:"size:100" json:"sector"`
	EntrepreneursMails  string           `json:"entrepreneurs_mails"`
	IsPastCampaign      bool             `gorm:"default:false" json:"is_past_campaign"`
	CampaignStatus      string           `gorm:"size:50;default:'pending'" json:"campaign_status"`
	CreatedAt           time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Prizes              []Prize          `gorm:"foreignKey:CampaignID;references:ID" json:"prizes"`
	Patents             []Patent         `gorm:"foreignKey:CampaignID;references:ID" json:"patents"`
	Law                 []Law            `gorm:"foreignKey:CampaignID;references:ID" json:"law"`
	TeamMembers         []TeamMember     `gorm:"foreignKey:CampaignID;references:ID" json:"team_members"`
	VisualVideos        []VisualVideo    `gorm:"foreignKey:CampaignID;references:ID" json:"visual_videos"`
	OtherDocuments      []OtherDocument  `gorm:"foreignKey:CampaignID;references:ID" json:"other_documents"`
	Risks               []Risk           `gorm:"foreignKey:CampaignID;references:ID" json:"risks"`
}

// Prize modeli (Kampanyaya bağlı ödüller)
type Prize struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID           uint      `gorm:"not null;index" json:"campaign_id"`
	PrizeDate            time.Time `json:"prize_date"`
	PrizeDescription     string    `json:"prize_description"`
	PrizePath            string    `json:"prize_path"`
	AwardingOrganization string    `gorm:"size:255" json:"awarding_organization"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Patent modeli (Kampanya ile ilişkili patent/dokümanlar)
type Patent struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID   uint      `gorm:"not null;index" json:"campaign_id"`
	DocumentNo   string    `gorm:"size:50" json:"document_no"`
	Description  string    `gorm:"type:text" json:"description"`
	DocumentPath string    `gorm:"type:text" json:"document_path"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Law modeli (Kampanya ile ilişkili hukuki izin dokümanları)
type Law struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID            uint      `gorm:"not null;index" json:"campaign_id"`
	PermissionSubject     string    `gorm:"size:255" json:"permission_subject"`     // Hukuki izin konusu
	PermissionPath        string    `gorm:"type:text" json:"permission_path"`         // İzin doküman yolu
	PermissionDescription string    `gorm:"type:text" json:"permission_description"`  // Açıklamalar
	CreatedAt             time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TeamMember modeli (Kampanya ile ilişkili ekip üyeleri)
type TeamMember struct {
	ID                     uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID             uint      `gorm:"not null;index" json:"campaign_id"`
	Photograph             string    `gorm:"type:text" json:"photograph"`              // Üye fotoğraf URL’si
	MembersName            string    `gorm:"size:100" json:"members_name"`             // Üye adı
	MembersSurname         string    `gorm:"size:100" json:"members_surname"`          // Üye soyadı
	MembersTitle           string    `gorm:"size:100" json:"members_title"`            // Ünvanı (CTO, CEO vb.)
	Resume                 string    `gorm:"type:text" json:"resume"`                  // CV bağlantısı
	Biography              string    `gorm:"type:text" json:"biography"`               // Kısa biyografi
	MembersTask            string    `gorm:"type:text" json:"members_task"`            // Üyenin görevleri
	MembersResponsibility  string    `gorm:"type:text" json:"members_responsibility"`  // Sorumlulukları
	EntrepreneurLinkMember string    `gorm:"type:text" json:"entrepreneur_link_member"`// Girişimci profili bağlantısı (Varsa)
	MembersMail            string    `gorm:"size:255" json:"members_mail"`             // E-posta adresi
	MembersInstagram       string    `gorm:"type:text" json:"members_instagram"`       // Instagram hesabı
	MembersTwitter         string    `gorm:"type:text" json:"members_twitter"`         // Twitter hesabı
	MembersLinkedin        string    `gorm:"type:text" json:"members_linkedin"`        // LinkedIn hesabı
	CreatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// VisualVideo modeli (Kampanya ile ilişkili tanıtım video ve görseller)
type VisualVideo struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID        uint      `gorm:"not null;index" json:"campaign_id"`
	DisplayPhotograph string    `gorm:"type:text" json:"display_photograph"`  // Kampanya ana görseli URL
	OtherPhotograph   string    `gorm:"type:text" json:"other_photograph"`    // Ek görseller URL
	VideoLink         string    `gorm:"type:text" json:"video_link"`          // Tanıtım videosu bağlantısı
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// OtherDocument modeli (Kampanyaya ait diğer belgeler)
type OtherDocument struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID uint      `gorm:"not null;index" json:"campaign_id"`
	FilePath   string    `gorm:"type:text" json:"file_path"` // Belgelerin saklandığı URL veya sistem yolu
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Risk modeli (Kampanya ile ilişkili risk bilgileri)
type Risk struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID      uint      `gorm:"not null;index" json:"campaign_id"`
	RiskType        string    `gorm:"size:50" json:"risk_type"`           // Örn: product, sector, partner, other
	RiskDescription string    `gorm:"type:text" json:"risk_description"`  // Risk açıklaması
	RiskMitigation  string    `gorm:"type:text" json:"risk_mitigation"`   // Riski azaltma planları
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
