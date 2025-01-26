# Backend Projesi (Go + Fiber)

Bu proje, Go dilinde Fiber framework kullanılarak geliştirilmiştir. Her **modül**, kendi rota dosyası (`*_routes.go`) dahil olmak üzere **controller**, **service**, **repository**, **model** ve **validation** bileşenlerini içerir. En üst seviyede, `router.go` dosyası tüm modül rotalarını birleştirerek uygulamaya kaydeder.

---

## 1. Nasıl Kurulur?

### 1.1. Gereksinimler

- [Go 1.22+](https://go.dev/dl/)
- PostgreSQL (veya tercihen bir başka veritabanı)
- (Opsiyonel) Docker ve Docker Compose
- (Opsiyonel) Migration aracı (golang-migrate, goose vb.)

### 1.2. Projeyi Çalıştırma

```bash
# Proje dizinine gir
cd backend/

# Modülleri indir
go mod download

# Migration (kullanılıyorsa)
# go run cmd/migrate.go up

# Projeyi derle ve çalıştır
go run cmd/server/main.go

## **3. Özet**

- **Her modül** (ör. `user`, `campaign`, `investment`) kendi rota dosyası (`*_routes.go`) ile kontrolü elinde tutar.
- **`router.go`** dosyası, tüm bu rota fonksiyonlarını (örn. `SetupUserRoutes`, `SetupCampaignRoutes`) tek bir yerden çağırarak Fiber’e kaydeder.
- Böylece **modülerlik** korunur, rota tanımları **dağınık** olmaz, ve ekip üyeleri **domain** bazında net bir şekilde ayrılmış kod üzerinde çalışabilir.


Proje Mimarisi & Dizini
Bu proje, katmanlı mimari ve modül bazlı bir yapı benimser. Her modül kendi rotasını (*_routes.go) barındırır, ardından internal/router.go tek bir noktadan bu rotaları yükler.

2.1. cmd/server/
Uygulamanın giriş noktası (main.go).
Fiber uygulamasını başlatır ve internal/server.go üzerinden konfigüre eder.
2.2. configs/
Konfigürasyon dosyaları (örn. config.yaml).
config.go: Bu dosyayı struct şeklinde uygulamaya aktaran kod.
2.3. internal/modules/
Burada her modül, kendi controller, service, repository, model, validation ve rota dosyalarına sahiptir.

Örnek: user/ modülü

user_controller.go    # HTTP handler'lar
user_service.go       # İş mantığı
user_repository.go    # DB erişimi
user_model.go         # GORM struct'lar (User, vb.)
user_validation.go    # Modüle özel doğrulama
user_routes.go        # Fiber route tanımları (func SetupUserRoutes())
user_test.go          # Modüle özel testler


*_routes.go dosyasında, o modüle ait endpoint’ler tanımlanır.
Örnek rota fonksiyonu


func SetupUserRoutes(router fiber.Router, ctrl *UserController) {
    router.Post("/users", ctrl.CreateUser)
    router.Get("/users/:id", ctrl.GetUser)
    ...
}

internal/middlewares/
Global veya proje genelinde kullanılan middleware fonksiyonları.
Örnek:
auth_middleware.go: JWT kontrolü, rol denetimi
logging_middleware.go: Log
rate_limit_middleware.go: Saldırı önleme


internal/router.go
Tüm modül rotalarını tek bir yerde birleştirir.
Örnek:
func SetupRoutes(app *fiber.App) {
    userCtrl := user.NewUserController(...)
    user.SetupUserRoutes(app, userCtrl)

    campaignCtrl := campaign.NewCampaignController(...)
    campaign.SetupCampaignRoutes(app, campaignCtrl)

    // Diğer modüllerin rota kayıtları ...
}


internal/server.go
Fiber uygulamasının yapılandırıldığı dosya.
InitServer() veya StartServer() fonksiyonları burada olur:

func InitServer() *fiber.App {
    app := fiber.New()
    // Middlewares
    // Logging
    // ...
    return app
}


internal/utils/
Ortak yardımcı fonksiyonlar (ör. JWT, şifreleme, hashing).
internal/validations/
Ortak validasyon fonksiyon veya struct’ları.
Modül bazlı validasyon yeterli değilse burayı kullanın.
migrations/
SQL tabanlı veya Go tabanlı migration dosyaları.
Veritabanı şemasını oluşturmak ve güncellemek için kullanın.
pkg/
Projeye dışarıdan eklenen veya bağımsız paket gibi kullanılan kodlar.
Örneğin logger veya emailer.
test/
Integration veya unit test koleksiyonları.
Örneğin Postman entegrasyon testleri integration/ altına koyulabilir.




Kod Yazım Standartları
Fonksiyon & Değişken İsimleri
Go konvansiyonlarına uyun (camelCase, PascalCase).
Dosya İsimleri
Genelde snake_case: user_controller.go, campaign_routes.go.
Rota Tanımları
Her modülün *_routes.go dosyasında SetupXxxRoutes() fonksiyonu bulunsun.
internal/router.go içerisinde bu fonksiyonlar çağrılarak rotalar toplu halde kayıt edilir.
Test Dosyaları
*_test.go formatı. Pozitif ve negatif senaryoları test edin.
Mümkünse testify veya benzeri kütüphaneler kullanın.
Error Yönetimi
Hataları logger paketi veya benzeri ile loglayın.
Kullanıcıya sade mesajlar döndürün (c.Status(fiber.StatusBadRequest).JSON(...)).
Commit Mesajları
Açıklayıcı, kısa özet + geniş açıklama.
PR incelemesinde mimari bütünlüğü koruyun.


Çalışma Akışı
Yeni Özellik Geliştirme
İlgili modül klasöründe (internal/modules/x_modul) çalışın.
Controller, Service, Repository, Model, Validation, Routes dosyalarını güncelleyin.
Test Süreci
Modül bazlı unit testler (*_test.go).
test/integration/ klasöründe API test (ör. Postman).
Kod İnceleme
Pull Request’lerde en az bir “deneyimli göz” onayı.
Proje standartlarına (mimari, isimlendirme vb.) uyum aranır.
Versiyonlama
Her sprint sonu bir tag (v0.1.0), major değişiklikte major versiyon artır.



Sık Karşılaşılan Sorunlar
CORS veya Rate Limit Sorunları
middlewares/ klasöründe doğru konfigüre edildiğinden emin olun.
Fiber Config
internal/server.go içinde body limit, zaman aşımı gibi ayarları yapabilirsiniz.
Veritabanı Bağlantısı
configs/config.go içindeki DB ayarlarına ve migration sırasına dikkat edin.


 İletişim ve Katkıda Bulunma
Pull Request: Testlerin sorunsuz geçmesi, en az bir review gereklidir.
Issue: Hata veya iyileştirme önerileri için GitHub Issues kullanın.
Sorumlular:
DevOps: Şefika https://github.com/searslan1
Backend: Seda, Ekrem, Emin, Berat, Gülay
Frontend: Mustafa
Product Owner: Erdal
(https://github.com/ekremsekmen)
(https://github.com/beroyzr)
(https://github.com/gulayyy)
(https://github.com/msfcodmsf)
(https://github.com/)
(https://github.com/)
(https://github.com/erdalgumuss)

 Lisans
Bu proje GAMMA TEAM Lisansı altında lisanslanmıştır. Detaylar için LICENSE dosyasına göz atabilirsiniz.



---

backend/
 ├── cmd/
 │    └── server/
 │         └── main.go
 ├── configs/
 │    ├── config.yaml
 │    └── config.go
 ├── internal/
 │    ├── modules/
 │    │    ├── user/
 │    │    │    ├── user_controller.go
 │    │    │    ├── user_service.go
 │    │    │    ├── user_repository.go
 │    │    │    ├── user_model.go
 │    │    │    ├── user_validation.go
 │    │    │    ├── user_routes.go         <-- Modüle özgü rota tanımları
 │    │    │    └── user_test.go
 │    │    ├── campaign/
 │    │    │    ├── campaign_controller.go
 │    │    │    ├── campaign_service.go
 │    │    │    ├── campaign_repository.go
 │    │    │    ├── campaign_model.go
 │    │    │    ├── campaign_validation.go
 │    │    │    ├── campaign_routes.go     <-- Modüle özgü rota tanımları
 │    │    │    └── campaign_test.go
 │    │    ├── investment/
 │    │    │    ├── investment_controller.go
 │    │    │    ├── investment_service.go
 │    │    │    ├── investment_repository.go
 │    │    │    ├── investment_model.go
 │    │    │    ├── investment_validation.go
 │    │    │    ├── investment_routes.go   <-- Modüle özgü rota tanımları
 │    │    │    └── investment_test.go
 │    │    ├── ... (diğer modüller) ...
 │    ├── middlewares/
 │    │    ├── auth_middleware.go
 │    │    ├── logging_middleware.go
 │    │    └── rate_limit_middleware.go
 │    ├── router.go                        <-- Tüm modül rotalarını birleştirir
 │    ├── server.go
 │    ├── utils/
 │    │    ├── crypto.go
 │    │    ├── jwt.go
 │    │    └── password.go
 │    └── validations/
 │         └── common_validations.go
 ├── migrations/
 │    ├── 20230101_init.sql
 │    └── 20230102_add_tables.sql
 ├── pkg/
 │    └── logger/
 │         ├── logger.go
 │         └── logger_test.go
 ├── test/
 │    ├── integration/
 │    │    └── postman_collection.json
 │    └── unit/
 │         └── ...
 ├── go.mod
 ├── go.sum
 └── README.md
```
