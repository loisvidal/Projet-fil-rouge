package database

import (
	"RedProject/models"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query("SELECT id, name, password, email, is_admin, is_confirmed FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.IdUser, &u.NameUser, &u.Password, &u.Mail, &u.IsAdmin, &u.IsConfirmed); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(id int) (models.User, error) {
	var u models.User
	err := DB.QueryRow("SELECT id, name, password, email, is_admin, is_confirmed FROM users WHERE id = ?", id).
		Scan(&u.IdUser, &u.NameUser, &u.Password, &u.Mail, &u.IsAdmin, &u.IsConfirmed)
	if err != nil {
		return u, err
	}
	return u, nil
}

func GetUserByLogin(identifier string) (models.User, error) {
	var u models.User
	err := DB.QueryRow(
		"SELECT id, name, password, email, is_admin, is_confirmed, failed_attempts, locked_until FROM users WHERE (name = ? OR email = ?)",
		identifier, identifier,
	).Scan(&u.IdUser, &u.NameUser, &u.Password, &u.Mail, &u.IsAdmin, &u.IsConfirmed, &u.FailedAttempts, &u.LockedUntil)
	if err != nil {
		return u, err
	}
	return u, nil
}

func CreateUser(u models.User) (int64, string, error) {
	token := generateToken()
	expiresAt := time.Now().Add(2 * time.Minute)

	hash, err := HashPassword(u.Password)
	if err != nil {
		return 0, "", fmt.Errorf("erreur hash: %v", err)
	}

	res, err := DB.Exec(
		"INSERT INTO users (name, password, email, confirmation_token, confirmation_expires_at) VALUES (?, ?, ?, ?, ?)",
		u.NameUser, hash, u.Mail, token, expiresAt,
	)
	if err != nil {
		return 0, "", err
	}
	id, _ := res.LastInsertId()
	return id, token, nil
}

func ConfirmUser(token string) error {
	res, err := DB.Exec(
		"UPDATE users SET is_confirmed = TRUE, confirmation_token = NULL, confirmation_expires_at = NULL WHERE confirmation_token = ? AND confirmation_expires_at > NOW()",
		token,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("lien de confirmation invalide ou expiré")
	}
	return nil
}

func IncrementFailedAttempts(userID int) error {
	_, err := DB.Exec(
		`UPDATE users SET failed_attempts = failed_attempts + 1,
		 locked_until = IF(failed_attempts >= 4, DATE_ADD(NOW(), INTERVAL 15 MINUTE), NULL)
		 WHERE id = ?`, userID)
	return err
}

func ResetFailedAttempts(userID int) error {
	_, err := DB.Exec("UPDATE users SET failed_attempts = 0, locked_until = NULL WHERE id = ?", userID)
	return err
}

func UpdatePassword(userID int, hash string) error {
	_, err := DB.Exec("UPDATE users SET password = ?, reset_token = NULL, reset_expires_at = NULL WHERE id = ?", hash, userID)
	return err
}

func UpdateUser(u models.User) error {
	_, err := DB.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?",
		u.NameUser, u.Mail, u.IdUser)
	return err
}

func DeleteUser(id int) error {
	_, err := DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func ToggleAdmin(id int) error {
	_, err := DB.Exec("UPDATE users SET is_admin = NOT is_admin WHERE id = ?", id)
	return err
}

func GetUserLikedProperties(userID int) ([]models.Property, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.name, p.description, p.price, p.is_sell, p.images, p.property_type, p.rooms, p.location, p.surface
		FROM properties p
		JOIN user_liked_properties ulp ON ulp.property_id = p.id
		WHERE ulp.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProperties(rows)
}

func GetUserBoughtProperties(userID int) ([]models.Property, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.name, p.description, p.price, p.is_sell, p.images, p.property_type, p.rooms, p.location, p.surface
		FROM properties p
		JOIN user_bought_properties ubp ON ubp.property_id = p.id
		WHERE ubp.user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProperties(rows)
}

func GetUserSellProperties(userID int) ([]models.Property, error) {
	rows, err := DB.Query(`
		SELECT id, name, description, price, is_sell, images, property_type, rooms, location, surface
		FROM properties
		WHERE owner_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProperties(rows)
}

func EnsureAdminAccount() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", "alexandre").Scan(&count)
	if count > 0 {
		return
	}

	log.Println("Création du compte admin alexandre...")
	hash, err := HashPassword("Ynov_123")
	if err != nil {
		log.Printf("Erreur hash: %v", err)
		return
	}
	_, err = DB.Exec(
		"INSERT INTO users (name, password, email, is_admin, is_confirmed) VALUES (?, ?, ?, TRUE, TRUE)",
		"alexandre", hash, "alexandre.petitfrere@ynov.com",
	)
	if err != nil {
		log.Printf("Erreur création admin: %v", err)
	}
}

func IsAdmin(userID int) (bool, error) {
	var admin bool
	err := DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", userID).Scan(&admin)
	if err != nil {
		return false, err
	}
	return admin, nil
}

func GetUserByIDWithRelations(userID int) (models.User, error) {
	u, err := GetUserByID(userID)
	if err != nil {
		return u, err
	}
	u.LikedProperty, _ = GetUserLikedProperties(userID)
	u.BuyProperty, _ = GetUserBoughtProperties(userID)
	u.SellProperty, _ = GetUserSellProperties(userID)
	return u, nil
}

func scanProperties(rows *sql.Rows) ([]models.Property, error) {
	var props []models.Property
	for rows.Next() {
		var p models.Property
		var imagesJSON sql.NullString
		var propType, location sql.NullString
		var rooms sql.NullInt64
		var surface sql.NullFloat64
		if err := rows.Scan(&p.IdProperty, &p.NameProperty, &p.DescProprety,
			&p.PriceProperty, &p.IsSellProperty, &imagesJSON,
			&propType, &rooms, &location, &surface); err != nil {
			return nil, err
		}
		if imagesJSON.Valid {
			p.ImgProperty = parseJSONArray(imagesJSON.String)
		}
		p.Type = propType.String
		p.Rooms = int(rooms.Int64)
		p.Location = location.String
		p.Surface = surface.Float64
		props = append(props, p)
	}
	return props, nil
}

func parseJSONArray(s string) []string {
	var arr []string
	if s == "" || s == "null" || s == "[]" {
		return arr
	}
	trimmed := s[1 : len(s)-1]
	if trimmed == "" {
		return arr
	}
	for i := 0; i < len(trimmed); i++ {
		if trimmed[i] == '"' {
			end := i + 1
			for end < len(trimmed) && trimmed[end] != '"' {
				if trimmed[end] == '\\' {
					end++
				}
				end++
			}
			arr = append(arr, trimmed[i+1:end])
			i = end
		}
	}
	return arr
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
