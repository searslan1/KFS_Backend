package utils

import (
    "os"
    "time"

    "github.com/sirupsen/logrus" // Loglama için popüler logrus paketi kullanılıyor
)

// Log değişkeni, uygulama genelinde kullanılacak logrus logger'ını temsil eder
var Log = logrus.New()

// RequestLog, bir HTTP isteğine ait loglanacak bilgileri tanımlayan yapı
type RequestLog struct {
    Method     string        // HTTP metodunu (GET, POST vb.) tutar
    Path       string        // İstek yapılan URL path'i
    IP         string        // İsteği yapan istemcinin IP adresi
    Duration   time.Duration // İstek işlenme süresi
    StatusCode int           // Yanıtın HTTP durum kodu
    UserAgent  string        // İstek yapan client'ın user-agent bilgisi
    ReqBody    string        // İstek gövdesi
    RespSize   int           // Yanıt boyutu (byte cinsinden)
    Error      error         // İşlem sırasında oluşan hata varsa ilgili hata bilgisi
    UserID     string        // İsteği yapan kullanıcıya ait ID (AuthMiddleware tarafından eklenir)
    Role       string        // İsteği yapan kullanıcının rolü (AuthMiddleware tarafından eklenir)
}

// init fonksiyonu, logger ilk başlatıldığında yapılandırmayı ayarlamak için kullanılır
func init() {
    // Logların JSON formatında yazılması için formatlayıcı ayarlanır.
    Log.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339, // Zaman damgası, RFC3339 formatında olacak
    })

    // Logları "api.log" adlı dosyaya ve stdout'a yazma
    // Dosya açılır ya da mevcut değilse oluşturulur; sürekli ekleme modunda açılır.
    file, err := os.OpenFile("api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err == nil {
        // Eğer dosya başarıyla açılırsa, log çıkışı bu dosyaya yönlendirilir.
        Log.SetOutput(file)
    }

    // Log seviyesini INFO olarak ayarla
    Log.SetLevel(logrus.InfoLevel)
}

// LogRequest fonksiyonu, gelen RequestLog yapısını alır ve loglama işlemini gerçekleştirir.
func LogRequest(req RequestLog) {
    // Loglanacak alanlar (fields) oluşturuluyor.
    fields := logrus.Fields{
        "method":      req.Method,                     // HTTP metodu
        "path":        req.Path,                       // İstek yapılan yol
        "ip":          req.IP,                         // İstemci IP adresi
        "duration_ms": req.Duration.Milliseconds(),    // İşlem süresi milisaniye cinsinden
        "status":      req.StatusCode,                 // Yanıtın durum kodu
        "user_agent":  req.UserAgent,                  // İstemci bilgisi (User-Agent)
        "req_body":    req.ReqBody,                    // İstek gövdesi
        "resp_size":   req.RespSize,                   // Yanıt boyutu
        "user_id":     req.UserID,                     // Kullanıcı ID'si
        "role":        req.Role,                       // Kullanıcı rolü
    }

    // Eğer istekte bir hata varsa, hata bilgisi eklenir ve ERROR seviyesinde loglanır.
    if req.Error != nil {
        fields["error"] = req.Error.Error() // Hata bilgisini string olarak ekliyoruz
        Log.WithFields(fields).Error("Request failed")
    } else {
        // Hata yoksa INFO seviyesinde başarılı log kaydı oluşturulur.
        Log.WithFields(fields).Info("Request completed")
    }
}