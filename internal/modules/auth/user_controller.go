package user

//seda
import (
	"encoding/json"
	"net/http"
)

type UserController struct {
	Service *UserService
}

// Kullanıcı kaydı endpoint'i
func (c *UserController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"user_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz giriş", http.StatusBadRequest)
		return
	}

	if err := c.Service.RegisterUser(req.Email, req.Password, req.UserType); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Kullanıcı başarıyla kaydedildi"))
}

// Kullanıcı giriş endpoint'i
func (c *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz giriş", http.StatusBadRequest)
		return
	}

	user, err := c.Service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
