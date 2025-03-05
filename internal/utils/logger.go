package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog" // Logrus yerine zerolog kullanıyoruz
)

// Log değişkeni, uygulama genelinde kullanılacak zerolog logger'ını temsil eder
var Log zerolog.Logger

// RequestLog, bir HTTP isteğine ait loglanacak bilgileri tanımlayan yapı
type RequestLog struct {
	Method      string        // HTTP metodunu (GET, POST vb.) tutar
	Path        string        // İstek yapılan URL path'i
	IP          string        // İsteği yapan istemcinin IP adresi
	Duration    time.Duration // İstek işlenme süresi
	StatusCode  int           // Yanıtın HTTP durum kodu
	UserAgent   string        // İstek yapan client'ın user-agent bilgisi
	ReqBodySize int           // İstek gövdesi boyutu (byte cinsinden)
	RespSize    int           // Yanıt boyutu (byte cinsinden)
	Error       error         // İşlem sırasında oluşan hata varsa ilgili hata bilgisi
	UserID      string        // İsteği yapan kullanıcıya ait ID (AuthMiddleware tarafından eklenir)
	Role        string        // İsteği yapan kullanıcının rolü (AuthMiddleware tarafından eklenir)
}

// init fonksiyonu, logger ilk başlatıldığında yapılandırmayı ayarlamak için kullanılır
func init() {
	// Zaman formatını ayarla
	zerolog.TimeFieldFormat = time.RFC3339

	// Logları "api.log" adlı dosyaya yaz
	file, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Dosya açılamazsa hata fırlat ve programı sonlandır
		panic(fmt.Sprintf("Log dosyası oluşturulamadı: %v", err))
	}

	// Logger'ı sadece dosyaya yazacak şekilde yapılandır
	Log = zerolog.New(file).
		With().
		Timestamp().
		Logger()

	// Log seviyesini INFO olarak ayarla
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// LogRequest fonksiyonu, gelen RequestLog yapısını alır ve loglama işlemini gerçekleştirir.
func LogRequest(req RequestLog) {
	// Log event oluştur
	var event *zerolog.Event

	// Hata varsa Error seviyesinde log al
	if req.Error != nil {
		event = Log.Error().Err(req.Error)
	} else {
		event = Log.Info()
	}

	// Log alanlarını ekle
	event.
		Str("method", req.Method).
		Str("path", req.Path).
		Str("ip", req.IP).
		Int64("duration_ms", req.Duration.Milliseconds()).
		Int("status", req.StatusCode).
		Str("user_agent", req.UserAgent).
		Int("req_size", req.ReqBodySize).
		Int("resp_size", req.RespSize).
		Str("user_id", req.UserID).
		Str("role", req.Role).
		Msg("HTTP Request") // Log mesajını gönder
}
