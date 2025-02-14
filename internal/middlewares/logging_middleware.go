package middlewares

import (
    "time"

    "KFS_Backend/internal/utils"

    "github.com/gofiber/fiber/v2"
)

func LoggingMiddleware(c *fiber.Ctx) error {
    // İsteğin başlangıç zamanını kaydet, bu süre isteğin işleme süresini hesaplamak için kullanılacak.
    start := time.Now()

    // Request bilgilerini topla:
    // İstek gövdesi (body) string olarak alınıyor.
    reqBody := string(c.Body())
    // Client'ın User-Agent başlığı alınıyor.
    userAgent := c.Get("User-Agent")

    // AuthMiddleware tarafından önceden context'e eklenmiş olan "userID" ve "role" bilgilerini çek.
    // Eğer bilgiler bulunamazsa, boş değer ile çalışmaya devam eder.
    userID, _ := c.Locals("userID").(string)
    userRole, _ := c.Locals("role").(string)

    // İstek işleniyor; c.Next() ile zincirdeki bir sonraki middleware veya handler'a geçiliyor.
    err := c.Next()

    // İstek işlendikten sonra yanıt bilgileri toplanıyor:
    // İşlemin süresi hesaplanıyor.
    duration := time.Since(start)
    // Yanıtın HTTP durum kodu alınıyor.
    statusCode := c.Response().StatusCode()
    // Yanıt boyutu (byte cinsinden) hesaplanıyor.
    respSize := len(c.Response().Body())

    // Toplanan tüm bilgileri içeren RequestLog yapısı oluşturuluyor
    // ve utils.LogRequest fonksiyonu ile loglama işlemi yapılıyor.
    utils.LogRequest(utils.RequestLog{
        Method:     c.Method(),    // HTTP metodu (GET, POST, vs.)
        Path:       c.Path(),      // İstek yapılan URL path
        IP:         c.IP(),        // İstemci IP adresi
        Duration:   duration,      // İstek işlenme süresi
        StatusCode: statusCode,    // Yanıtın durum kodu
        UserAgent:  userAgent,     // İstemcinin User-Agent bilgisi
        ReqBody:    reqBody,       // İstek gövdesi
        RespSize:   respSize,      // Yanıt boyutu (byte cinsinden)
        Error:      err,           // İşlem sırasında meydana gelen hata (varsa)
        UserID:     userID,        // İstemciye ait kullanıcı ID'si
        Role:       userRole,      // İstemcinin rolü
    })

    // Eğer işleme sırasında bir hata oluştuysa, bu hata sonraki middleware ya da handler'a iletiliyor.
    return err
}

/*
LoggingMiddleware Açıklaması:
- Bu middleware, gelen HTTP isteklerinin detaylı loglarını toplar.
- İstek başladığında zaman damgası alınır; istek gövdesi, User-Agent, IP, URL path gibi bilgiler toplanır.
- AuthMiddleware tarafından daha önce eklenmiş olan kullanıcı kimlik bilgileri (userID, role) de burada elde edilir.
- c.Next() ile zincirdeki sonraki middleware veya handler çalıştırıldıktan sonra, istek işlenme süresi, yanıtın durum kodu ve yanıt boyutu hesaplanır.
- Tüm bu veriler, utils.LogRequest fonksiyonu aracılığıyla loglanır.
Kullanım:
- Fiber uygulamanızda, bu middleware'i global olarak ekleyebilirsiniz; örneğin main.go dosyanızda app.Use(LoggingMiddleware) şeklinde.
- Böylece, uygulamaya gelen her istek için detaylı log kayıtları oluşturulur ve sistem izlenebilir hale gelir.
*/