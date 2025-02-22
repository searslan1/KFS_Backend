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
	CampaignName        string           `gorm:"size:255" json:"campaign_name"`
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
	ProductProduction   []ProductProduction `gorm:"foreignKey:CampaignID;references:ID" json:"product_production"`
	Risks               []Risks          `gorm:"foreignKey:CampaignID;references:ID" json:"risks"`
	Funding 		    []Funding        `gorm:"foreignKey:CampaignID;references:ID" json:"funding"`
	Establishment		[]Establishment  `gorm:"foreignKey:CampaignID;references:ID" json:"establishment"`
	FinancialTables     []FinancialTable `gorm:"foreignKey:CampaignID;references:ID" json:"financial_tables"`
	VisualVideos        []VisualVideo    `gorm:"foreignKey:CampaignID;references:ID" json:"visual_videos"`
	OtherDocuments      []OtherDocuments `gorm:"foreignKey:CampaignID;references:ID" json:"other_documents"`
	
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

// Patent modeli (Kampanya ile ilişkili patent/belge)
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

// TeamMember modeli (Kampanya ile ilişkili ekip üyeleri Takım)
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

// ProductProductionModel modeli
type ProductProduction struct {
	ID                                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID                          uint      `gorm:"not null;index" json:"campaign_id"` // Kampanyaya bağlı
	ProductSummary                      string    `gorm:"type:text" json:"product_summary"`  // Ürün özeti
	ProductAbout                        string    `gorm:"type:text" json:"product_about"`    // Ürün hakkında bilgi
	ProductPath                         string    `gorm:"type:text" json:"product_path"`     // Ürüne dair doküman
	ProductProblem                      string    `gorm:"type:text" json:"product_problem"`  // Çözülmesi hedeflenen problem
	ProductSolution                     string    `gorm:"type:text" json:"product_solution"` // Ürünün sunduğu çözüm
	ProductEvaluation                   string    `gorm:"type:text" json:"product_evaluation"` // Ürün değerlendirme raporu
	ProductsDevelopmentStageSummary     string    `gorm:"type:text" json:"products_development_stage_summary"`      // Ürün geliştirme aşaması
	ProductsDevelopmentStageSummaryPath string    `gorm:"type:text" json:"products_development_stage_summary_path"` // İlgili doküman yolu
	ProductsProductionStageSummary      string    `gorm:"type:text" json:"products_production_stage_summary"`       // Üretim süreci özeti
	ProductsProductionStageSummaryPath  string    `gorm:"type:text" json:"products_production_stage_summary_path"`  // Üretim süreci dosya yolu
	SideProductSummary                  string    `gorm:"type:text" json:"side_product_summary"`  // Yan ürün bilgisi
	SideProductSummaryPath              string    `gorm:"type:text" json:"side_product_summary_path"` // Yan ürün raporu
	AnalysisSummary                     string    `gorm:"type:text" json:"analysis_summary"`      // Analiz özeti
	AnalysisSummaryPath                 string    `gorm:"type:text" json:"analysis_summary_path"` // Analiz dokümanı yolu
	ARGESummary                         string    `gorm:"type:text" json:"ar_ge_summary"`         // AR-GE çalışmaları özeti
	ARGESummaryPath                     string    `gorm:"type:text" json:"ar_ge_summary_path"`    // AR-GE rapor dokümanı yolu
	PreviewSalesSummary                 string    `gorm:"type:text" json:"preview_sales_summary"` // Satış tahminleri özeti
	PreviewSalesSummaryPath             string    `gorm:"type:text" json:"preview_sales_summary_path"` // Satış öngörüleri dokümanı
	CreatedAt                           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// OtherSubject modeli (Ürün üretim süreciyle bağlantılı)
type OtherSubject struct {
	ID                        uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductProductionModelsID uint             `gorm:"not null;index" json:"product_production_models_id"` // Ürün üretim süreciyle bağlantılı
	Subject                   string           `gorm:"type:text" json:"subject"`                           // Konu başlığı
	Path                      string           `gorm:"type:text" json:"path"`                              // İlgili doküman yolu
	Description               string           `gorm:"type:text" json:"description"`                       // Açıklama
	CreatedAt                 time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                 time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Foreign key ilişkilendirmesi: ProductProductionModelsID sütunu,
	// ProductProduction (product_production_models tablosu) modelindeki ID sütununa bağlıdır.
	ProductProduction         ProductProduction `gorm:"foreignKey:ProductProductionModelsID;references:ID" json:"product_production"`
}

//Pazar/Rekabet/Hedef
//Analiz

// FinancialTable modeli (Kampanyaya bağlı finansal veriler)
type FinancialTable struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID          uint      `gorm:"not null;index" json:"campaign_id"`                // Kampanyaya bağlı finansal veriler
	InvestmentYear1     float64   `gorm:"type:decimal(15,2)" json:"investment_year_1"`       // 1. yıl yatırımı
	InvestmentYear2     float64   `gorm:"type:decimal(15,2)" json:"investment_year_2"`       // 2. yıl yatırımı
	InvestmentYear3     float64   `gorm:"type:decimal(15,2)" json:"investment_year_3"`       // 3. yıl yatırımı
	InvestmentYear4     float64   `gorm:"type:decimal(15,2)" json:"investment_year_4"`       // 4. yıl yatırımı
	InvestmentYear5     float64   `gorm:"type:decimal(15,2)" json:"investment_year_5"`       // 5. yıl yatırımı
	Total               float64   `gorm:"type:decimal(15,2)" json:"total"`                  // Toplam yatırım (year_1 + year_2 + ...)
	ProfitEstimatesText string    `gorm:"type:text" json:"profit_estimates_text"`           // Kar tahmini açıklaması
	ExplanationText     string    `gorm:"type:text" json:"explanation_text"`                // Ek açıklamalar
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Budget modeli (Finans tablosuyla ilişkili bütçe bilgileri)
type Budget struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FinancialTableID uint           `gorm:"not null;index" json:"financial_table_id"` // Finans tablosuyla ilişkili
	Category         string         `gorm:"size:100" json:"category"`                 // Bütçe kategorisi (Pazarlama, Ar-Ge, Üretim vb.)
	SubCategory      string         `gorm:"size:100" json:"sub_category"`             // Alt kategori
	Year1            float64        `gorm:"type:decimal(15,2)" json:"year_1"`         // 1. yıl bütçesi
	Year2            float64        `gorm:"type:decimal(15,2)" json:"year_2"`         // 2. yıl bütçesi
	Year3            float64        `gorm:"type:decimal(15,2)" json:"year_3"`         // 3. yıl bütçesi
	Year4            float64        `gorm:"type:decimal(15,2)" json:"year_4"`         // 4. yıl bütçesi
	Year5            float64        `gorm:"type:decimal(15,2)" json:"year_5"`         // 5. yıl bütçesi
	Total            float64        `gorm:"type:decimal(15,2)" json:"total"`          // Toplam bütçe (year_1 + year_2 + ...)
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Foreign key ilişkilendirmesi: FinancialTableID, FinancialTable modelinin ID sütununa bağlıdır.
	FinancialTable FinancialTable `gorm:"foreignKey:FinancialTableID;references:ID" json:"financial_table"`
}

// Document modeli (Finans tablosuyla ilişkili doküman bilgileri)
type Document struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FinancialTableID uint           `gorm:"not null;index" json:"financial_table_id"` // Finans tablosuyla ilişkili
	DocumentName     string         `gorm:"size:255" json:"document_name"`            // Belge adı
	DocumentURL      string         `gorm:"type:text" json:"document_url"`            // Belge bağlantısı
	SecureToken      string         `gorm:"type:text" json:"secure_token"`            // Doküman URL’leri için güvenlik artırılmalı. Çözüm: ALTER TABLE document ADD COLUMN secure_token TEXT;
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Foreign key ilişkilendirmesi: FinancialTableID FinancialTable modelinin ID sütununa bağlıdır.
	FinancialTable FinancialTable `gorm:"foreignKey:FinancialTableID;references:ID" json:"financial_table"`
}

// ProductRevenue modeli (Finansal tablo ile ilişkili ürün gelir bilgileri)
type ProductRevenue struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FinancialTableID uint           `gorm:"not null;index" json:"financial_table_id"` // Finansal tablo ile ilişkili
	ProductName      string         `gorm:"size:255" json:"product_name"`             // Ürün adı
	SalesPrice       float64        `gorm:"type:decimal(15,2)" json:"sales_price"`      // Satış fiyatı
	DirectCost       float64        `gorm:"type:decimal(15,2)" json:"direct_cost"`      // Üretim maliyeti
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Foreign key ilişkilendirmesi: FinancialTableID, FinancialTable modelinin ID sütununa bağlıdır.
	FinancialTable FinancialTable `gorm:"foreignKey:FinancialTableID;references:ID" json:"financial_table"`
}

// SalesTargets modeli (Ürün gelir tablosuyla bağlantılı satış hedefleri)
type SalesTargets struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductRevenueID uint           `gorm:"not null;index" json:"product_revenue_id"` // Ürün gelir tablosuyla bağlantılı
	TargetType       string         `gorm:"size:100" json:"target_type"`              // Hedef türü (adetsel, cirosal vb.)
	Year1            float64        `gorm:"type:decimal(15,2)" json:"year_1"`         // 1. yıl satış hedefi
	Year2            float64        `gorm:"type:decimal(15,2)" json:"year_2"`         // 2. yıl satış hedefi
	Year3            float64        `gorm:"type:decimal(15,2)" json:"year_3"`         // 3. yıl satış hedefi
	Year4            float64        `gorm:"type:decimal(15,2)" json:"year_4"`         // 4. yıl satış hedefi
	Year5            float64        `gorm:"type:decimal(15,2)" json:"year_5"`         // 5. yıl satış hedefi
	Total            float64        `gorm:"type:decimal(15,2)" json:"total"`          // Toplam satış hedefi (year_1 + year_2 + ...)
	CreatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Foreign key ilişkilendirmesi: ProductRevenueID, ProductRevenue modelinin ID sütununa bağlıdır.
	ProductRevenue   ProductRevenue `gorm:"foreignKey:ProductRevenueID;references:ID" json:"product_revenue"`
}

// Funding modeli (Kampanya ile bağlantılı finansman bilgileri)
type Funding struct {
	ID                           uint                      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID                   uint                      `gorm:"not null;index" json:"campaign_id"`                      // Kampanya ile bağlantılı
	InterferenceValue            float64                   `gorm:"type:decimal(15,2)" json:"interference_value"`           // Girişim için harcama maliyeti
	UseTime                      string                    `gorm:"type:text" json:"use_time"`                              // Fonun kullanım süresi
	EvaluationReport             string                    `gorm:"type:text" json:"evaluation_report"`                     // Değerleme raporu
	NeedAmountFund               float64                   `gorm:"type:decimal(15,2)" json:"need_amount_fund"`              // Gerekli yatırım miktarı
	AmountGivenShare             float64                   `gorm:"type:decimal(15,2)" json:"amount_given_share"`            // Verilen hisse miktarı
	NumberSaleShare              int                       `json:"number_sale_share"`                                      // Satılan hisse sayısı
	PostFundingCapital           float64                   `gorm:"type:decimal(15,2)" json:"post_funding_capital"`          // Finansman sonrası sermaye
	SalesSharePrice              float64                   `gorm:"type:decimal(15,2)" json:"sales_share_price"`             // Hisse satış fiyatı
	SalesNominalSharePrice       float64                   `gorm:"type:decimal(15,2)" json:"sales_nominal_share_price"`     // Hisse nominal fiyatı
	PlaceFundsCollectedID        uint                      `gorm:"not null;index" json:"place_funds_collected_id"`          // Toplanan fonlar bağlantısı
	AdditionalSourcesFinancingID uint                      `gorm:"not null;index" json:"additional_sources_financing_id"`   // Ek fon kaynakları bağlantısı
	ComparisonCurrentPostfunding string                    `gorm:"type:text" json:"comparison_current_postfunding"`         // Finansman öncesi-sonrası kıyaslama
	BasicInformation             string                    `gorm:"type:text" json:"basic_information"`                     // Ek açıklamalar
	CreatedAt                    time.Time                 `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt                    time.Time                 `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Foreign key ilişkileri
	PlaceFundsCollected          PlaceFundsCollected       `gorm:"foreignKey:PlaceFundsCollectedID;references:ID" json:"place_funds_collected"`
	AdditionalSourcesFinancing   AdditionalSourcesFinancing `gorm:"foreignKey:AdditionalSourcesFinancingID;references:ID" json:"additional_sources_financing"`
}

// AdditionalSourcesFinancing modeli (Ek fon kaynakları)
type AdditionalSourcesFinancing struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProvisionTime time.Time `json:"provision_time"`               // Fon sağlanma zamanı
	Description   string    `gorm:"type:text" json:"description"` // Açıklama
	Amount        float64   `gorm:"type:decimal(15,2)" json:"amount"` // Fon miktarı
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// PlaceFundsCollected modeli (Toplanan fonlar bilgileri)
type PlaceFundsCollected struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UsageStartTime  time.Time `json:"usage_start_time"`  // Fonun kullanılmaya başlandığı tarih
	UsageFinishTime time.Time `json:"usage_finish_time"` // Fon kullanımının tamamlandığı tarih
	Description     string    `gorm:"type:text" json:"description"` // Fon kullanım açıklaması
	Amount          float64   `gorm:"type:decimal(15,2)" json:"amount"` // Kullanılan fon miktarı
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Risk modeli (Kampanya ile ilişkili risk bilgileri)
type Risks struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID      uint      `gorm:"not null;index" json:"campaign_id"`
	RiskType        string    `gorm:"size:50" json:"risk_type"`           // Örn: product, sector, partner, other
	RiskDescription string    `gorm:"type:text" json:"risk_description"`  // Risk açıklaması
	RiskMitigation  string    `gorm:"type:text" json:"risk_mitigation"`   // Riski azaltma planları
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


// PostfundingPartners modeli (Finansman sonrası ortak bilgileri)
type PostfundingPartners struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                  string    `gorm:"size:255" json:"name"`                          // Ortak adı
	Task                  string    `gorm:"type:text" json:"task"`                         // Ortağın görevi
	University            string    `gorm:"size:255" json:"university"`                    // Akademik geçmişi (varsa)
	AverageNot            float64   `gorm:"type:decimal(3,2)" json:"average_not"`          // Akademik not ortalaması (varsa)
	Resume                string    `gorm:"type:text" json:"resume"`                       // CV bilgisi
	Citizenship           string    `gorm:"size:50" json:"citizenship"`                    // Vatandaşlık durumu
	ShareInCapitalAmount  float64   `gorm:"type:decimal(15,2)" json:"share_in_capital_amount"` // Şirkette sahip olduğu sermaye
	ShareInCapitalRate    float64   `gorm:"type:decimal(5,2)" json:"share_in_capital_rate"`    // Şirketteki hisse oranı (%)
	Vote                  bool      `json:"vote"`                                          // Oy hakkı var mı?
	Concession            string    `gorm:"type:text" json:"concession"`                   // Ayrıcalıklı haklar
	CampaignRelation      string    `gorm:"type:text" json:"campaign_relation"`            // Kampanya ile ilişki durumu
	WorkExperience        string    `gorm:"type:text" json:"work_experience"`              // Çalışma geçmişi
	AreasOfSpecialization string    `gorm:"type:text" json:"areas_of_specialization"`      // Uzmanlık alanları
	CreatedAt             time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Establishment modeli (Kampanya ile bağlantılı şirket kuruluş bilgileri)
type Establishment struct {
	ID                    uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID            uint                `gorm:"not null;index" json:"campaign_id"`           // Kampanya bağlantısı
	Title                 string              `gorm:"size:255" json:"title"`                       // Şirketin ismi
	Capitalization        float64             `gorm:"type:decimal(15,2)" json:"capitalization"`    // Şirketin ilk sermayesi
	City                  string              `gorm:"size:100" json:"city"`                        // Şirketin bulunduğu şehir
	District              string              `gorm:"size:100" json:"district"`                    // Şirketin bulunduğu ilçe
	Address               string              `gorm:"type:text" json:"address"`                    // Şirket adresi
	PostfundingPartnersID uint                `gorm:"not null;index" json:"postfunding_partners_id"` // Finansman sonrası ortaklar bağlantısı
	CreatedAt             time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt             time.Time           `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Foreign key ilişkilendirmesi: PostfundingPartnersID, PostfundingPartners modelinin ID sütununa bağlıdır.
	PostfundingPartners   PostfundingPartners `gorm:"foreignKey:PostfundingPartnersID;references:ID" json:"postfunding_partners"`
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
type OtherDocuments struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CampaignID uint      `gorm:"not null;index" json:"campaign_id"`
	FilePath   string    `gorm:"type:text" json:"file_path"` // Belgelerin saklandığı URL veya sistem yolu
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


