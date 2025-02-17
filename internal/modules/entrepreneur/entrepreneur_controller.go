package controller

import (
	"entrepreneur/model"
	"entrepreneur/service"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// EntrepreneurController girişimci işlemlerini yönetir
type EntrepreneurController struct {
	EntrepreneurService *service.EntrepreneurService
}

// **1️⃣ Girişimci Profili Oluşturma**
func (ec *EntrepreneurController) CreateEntrepreneur(c *fiber.Ctx) error {
	// Kullanıcı ID'sini al (Locals içinde olduğundan, öncelikle tür kontrolü yap)
	userIDInterface := c.Locals("userID")
	userID, ok := userIDInterface.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Kimlik doğrulama hatası.",
		})
	}

	// **SPK Gerekliliği: E-Devlet onayı kontrolü**
	if !ec.EntrepreneurService.IsEDevletVerified(userID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Kullanıcı e-Devlet doğrulamasından geçmemiş.",
		})
	}

	// **Kullanıcının zaten girişimci profili var mı?**
	existingProfile, err := ec.EntrepreneurService.GetByUserID(userID)
	if err == nil && existingProfile != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Kullanıcı zaten bir girişimci profiline sahip.",
		})
	}

	// **Request Body'yi parse et**
	var entrepreneur model.Entrepreneur
	if err := c.BodyParser(&entrepreneur); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz girişimci verisi. Lütfen girişimci profili bilgilerini kontrol edin.",
		})
	}

	// **Girişimci profilini oluştur**
	entrepreneur.UserID = userID
	entrepreneur.Status = "pending" // Admin onayı bekleniyor

	if err := ec.EntrepreneurService.Create(&entrepreneur); err != nil {
		log.Println("Girişimci profili oluşturulamadı:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Girişimci profili oluşturulamadı.",
		})
	}

	// **Başarılı dönüş**
	return c.Status(fiber.StatusCreated).JSON(entrepreneur)
}

// **2️⃣ Admin Girişimci Onayı**
func (ec *EntrepreneurController) AdminApproveEntrepreneur(c *fiber.Ctx) error {
	entrepreneurID := c.Params("id")

	// **ID'nin integer olup olmadığını kontrol et**
	entrepreneurIDParsed, err := strconv.Atoi(entrepreneurID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Geçersiz girişimci ID'si.",
		})
	}

	// **Girişimci profili veritabanından alınıyor**
	entrepreneur, err := ec.EntrepreneurService.GetByID(uint(entrepreneurIDParsed)) // ❌ HATA: GetByUserID yerine GetByID kullanmalısın
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Girişimci profili bulunamadı.",
		})
	}

	// **Girişimci profili zaten onaylanmış mı?**
	if entrepreneur.Status == "approved" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bu girişimci zaten onaylanmış.",
		})
	}

	// **Girişimciyi onayla**
	entrepreneur.Status = "approved"
	entrepreneur.IsAdminApproved = true

	// **Girişimci profili güncelleniyor**
	if err := ec.EntrepreneurService.Update(entrepreneur); err != nil {
		log.Println("Girişimci profili güncellenemedi:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Girişimci profili güncellenemedi.",
		})
	}

	// **Başarılı dönüş**
	return c.JSON(fiber.Map{
		"message": "Girişimci profili başarıyla onaylandı.",
		"data":    entrepreneur,
	})
}
