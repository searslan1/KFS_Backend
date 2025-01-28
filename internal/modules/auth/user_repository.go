package user

import (
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

// Kullanıcı ekleme
func (repo *UserRepository) CreateUser(user User) error {
	query := "INSERT INTO users (email, password_hash, user_type, created_at) VALUES (?, ?, ?, ?)"
	_, err := repo.DB.Exec(query, user.Email, user.Password, user.UserType, user.CreatedAt)
	return err
}

// Kullanıcıyı ID'ye göre getir
func (repo *UserRepository) GetUserByID(id int) (*User, error) {
	query := "SELECT id, email, password_hash, user_type, created_at FROM users WHERE id = ?"
	row := repo.DB.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.UserType, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (repo *UserRepository) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, email, password_hash, user_type, created_at FROM users WHERE email = ?"
	row := repo.DB.QueryRow(query, email)

	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.UserType, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Kullanıcı bulunamadıysa hata yerine `nil` döndür
		}
		return nil, err // Başka bir hata varsa döndür
	}
	return user, nil
}