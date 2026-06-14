package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitDB() {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASS", "")
	name := getEnv("DB_NAME", "yplaza")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4", user, pass, host, port, name)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erreur connexion MySQL: %v", err)
	}
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	if err = DB.Ping(); err != nil {
		log.Fatalf("Impossible de ping MySQL: %v", err)
	}

	log.Println("Connecté à MySQL (yplaza)")
	migrate()
}

func migrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE,
			is_confirmed BOOLEAN NOT NULL DEFAULT FALSE,
			confirmation_token VARCHAR(64) DEFAULT NULL,
			confirmation_expires_at DATETIME DEFAULT NULL,
			reset_token VARCHAR(64) DEFAULT NULL,
			reset_expires_at DATETIME DEFAULT NULL,
			failed_attempts INT DEFAULT 0,
			locked_until DATETIME DEFAULT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS properties (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10,2) NOT NULL DEFAULT 0,
			is_sell BOOLEAN NOT NULL DEFAULT TRUE,
			images JSON,
			owner_id INT,
			property_type VARCHAR(50) DEFAULT 'house',
			rooms INT DEFAULT 0,
			location VARCHAR(255) DEFAULT '',
			surface DECIMAL(8,2) DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS user_liked_properties (
			user_id INT NOT NULL,
			property_id INT NOT NULL,
			PRIMARY KEY (user_id, property_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS user_bought_properties (
			user_id INT NOT NULL,
			property_id INT NOT NULL,
			bought_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, property_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS auctions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			property_id INT NOT NULL,
			seller_id INT NOT NULL,
			start_price DECIMAL(10,2) NOT NULL DEFAULT 0,
			current_price DECIMAL(10,2) NOT NULL DEFAULT 0,
			min_bid_step DECIMAL(10,2) NOT NULL DEFAULT 100,
			winner_id INT DEFAULT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (property_id) REFERENCES properties(id) ON DELETE CASCADE,
			FOREIGN KEY (seller_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (winner_id) REFERENCES users(id) ON DELETE SET NULL
		)`,
		`CREATE TABLE IF NOT EXISTS bids (
			id INT AUTO_INCREMENT PRIMARY KEY,
			auction_id INT NOT NULL,
			user_id INT NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (auction_id) REFERENCES auctions(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS sales_history (
			id INT AUTO_INCREMENT PRIMARY KEY,
			property_id INT NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			property_type VARCHAR(50) DEFAULT 'house',
			rooms INT DEFAULT 0,
			location VARCHAR(255) DEFAULT '',
			surface DECIMAL(8,2) DEFAULT 0,
			sold_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Erreur migration: %v", err)
		}
	}

	alterQueries := []string{
		"ALTER TABLE properties ADD COLUMN property_type VARCHAR(50) DEFAULT 'house' AFTER price",
		"ALTER TABLE properties ADD COLUMN rooms INT DEFAULT 0 AFTER property_type",
		"ALTER TABLE properties ADD COLUMN location VARCHAR(255) DEFAULT '' AFTER rooms",
		"ALTER TABLE properties ADD COLUMN surface DECIMAL(8,2) DEFAULT 0 AFTER location",
		"ALTER TABLE users ADD COLUMN is_confirmed BOOLEAN NOT NULL DEFAULT FALSE AFTER is_admin",
		"ALTER TABLE users ADD COLUMN confirmation_token VARCHAR(64) DEFAULT NULL AFTER is_confirmed",
		"ALTER TABLE users ADD COLUMN confirmation_expires_at DATETIME DEFAULT NULL AFTER confirmation_token",
		"ALTER TABLE users ADD COLUMN reset_token VARCHAR(64) DEFAULT NULL AFTER confirmation_expires_at",
		"ALTER TABLE users ADD COLUMN reset_expires_at DATETIME DEFAULT NULL AFTER reset_token",
		"ALTER TABLE users ADD COLUMN failed_attempts INT DEFAULT 0 AFTER reset_expires_at",
		"ALTER TABLE users ADD COLUMN locked_until DATETIME DEFAULT NULL AFTER failed_attempts",
	}
	for _, q := range alterQueries {
		DB.Exec(q)
	}

	hashAdminPassword()
	log.Println("Tables migrées avec succès")
}

func hashAdminPassword() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM users WHERE name = 'alexandre' AND password NOT LIKE '$2a$%'").Scan(&count)
	if count == 0 {
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("Ynov_123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Erreur hash admin: %v", err)
		return
	}
	DB.Exec("UPDATE users SET password = ? WHERE name = 'alexandre'", string(hash))
	log.Println("Mot de passe admin hashé (bcrypt)")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
