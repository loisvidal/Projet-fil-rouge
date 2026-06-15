package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func ResetAndSeed() {
	log.Println("Réinitialisation de la base de données...")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASS", "")

	tempDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&charset=utf8mb4", user, pass, host, port)
	tempDB, err := sql.Open("mysql", tempDSN)
	if err != nil {
		log.Fatalf("Erreur connexion MySQL (seed): %v", err)
	}

	tempDB.Exec("DROP DATABASE IF EXISTS yplaza")
	_, err = tempDB.Exec("CREATE DATABASE yplaza CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	if err != nil {
		log.Fatalf("Erreur création base yplaza: %v", err)
	}
	tempDB.Close()

	log.Println("Base yplaza créée, exécution migrations...")
	InitDB()

	log.Println("Insertion des utilisateurs...")
	seedUsers()

	log.Println("Insertion des propriétés...")
	seedProperties()

	log.Println("Insertion des enchères...")
	seedAuctions()

	log.Println("Insertion des historiques de vente...")
	seedSalesHistory()

	log.Println("Insertion des relations utilisateur...")
	seedUserRelations()

	log.Println("Base de données initialisée avec succès !")
}

func hashBcrypt(password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erreur hash bcrypt: %v", err)
	}
	return string(h)
}

func seedUsers() {
	users := []struct {
		Name     string
		Password string
		Email    string
		IsAdmin  bool
	}{
		{"toto", hashBcrypt("titi123"), "toto@yopmail.com", false},
		{"Alexandre", hashBcrypt("En_78270"), "alexandre.petitfrere14@gmail.com", true},
	}
	for _, u := range users {
		_, err := DB.Exec(
			"INSERT INTO users (name, password, email, is_admin, is_confirmed) VALUES (?, ?, ?, ?, TRUE)",
			u.Name, u.Password, u.Email, u.IsAdmin,
		)
		if err != nil {
			log.Fatalf("Erreur insertion user %s: %v", u.Name, err)
		}
		log.Printf("  Utilisateur créé: %s (admin=%v)", u.Name, u.IsAdmin)
	}
}

func seedProperties() {
	now := time.Now()
	props := []struct {
		Name        string
		Desc        string
		Price       float64
		OwnerID     int
		PropType    string
		Rooms       int
		Location    string
		Surface     float64
	}{
		{"Villa Moderne avec Piscine", "Superbe villa contemporaine avec piscine à débordement, vue panoramique et jardin paysager.", 650000, 2, "villa", 5, "Nice", 180},
		{"Appartement Centre-Ville", "Bel appartement lumineux en plein centre-ville, proche commerces et transports.", 280000, 2, "apartment", 3, "Lyon", 72},
		{"Studio Lumineux", "Studio idéal pour étudiant, entièrement rénové avec kitchenette équipée.", 95000, 1, "studio", 1, "Toulouse", 25},
		{"Maison de Campagne", "Charmante maison de campagne avec terrain arboré, calme et vue dégagée.", 195000, 1, "house", 4, "Aix-en-Provence", 120},
		{"Loft Industriel", "Loft spacieux style new-yorkais avec hauts plafonds, poutres apparentes et mezzanine.", 340000, 2, "loft", 2, "Paris", 150},
		{"Local Commercial", "Local commercial de 200m² en plein centre, idéal pour restauration ou commerce.", 480000, 2, "commercial", 3, "Marseille", 200},
		{"Penthouse Vue Mer", "Superbe penthouse avec terrasse panoramique et vue imprenable sur la mer.", 890000, 1, "penthouse", 4, "Cannes", 160},
		{"Maison de Ville Traditionnelle", "Magnifique maison de ville rénovée, alliant charme ancien et confort moderne.", 420000, 2, "townhouse", 4, "Bordeaux", 140},
	}
	for _, p := range props {
		_, err := DB.Exec(
			`INSERT INTO properties (name, description, price, is_sell, owner_id, property_type, rooms, location, surface, created_at)
			 VALUES (?, ?, ?, TRUE, ?, ?, ?, ?, ?, ?)`,
			p.Name, p.Desc, p.Price, p.OwnerID, p.PropType, p.Rooms, p.Location, p.Surface, now,
		)
		if err != nil {
			log.Fatalf("Erreur insertion property %s: %v", p.Name, err)
		}
		log.Printf("  Propriété créée: %s (%.0f€, %s)", p.Name, p.Price, p.PropType)
	}
}

func seedAuctions() {
	now := time.Now()
	auctions := []struct {
		PropertyID int
		SellerID   int
		StartPrice float64
		MinStep    float64
		EndHours   int
	}{
		{1, 2, 550000, 5000, 72},
		{3, 1, 80000, 1000, 48},
		{5, 2, 300000, 3000, 96},
	}
	for _, a := range auctions {
		endTime := now.Add(time.Duration(a.EndHours) * time.Hour)
		_, err := DB.Exec(
			`INSERT INTO auctions (property_id, seller_id, start_price, current_price, min_bid_step, start_time, end_time, is_active)
			 VALUES (?, ?, ?, ?, ?, ?, ?, TRUE)`,
			a.PropertyID, a.SellerID, a.StartPrice, a.StartPrice, a.MinStep, now, endTime,
		)
		if err != nil {
			log.Fatalf("Erreur insertion auction prop %d: %v", a.PropertyID, err)
		}
		log.Printf("  Enchère créée: propriété #%d (%.0f€, %dh)", a.PropertyID, a.StartPrice, a.EndHours)
	}

	bids := []struct {
		AuctionID int
		UserID    int
		Amount    float64
	}{
		{1, 1, 560000},
		{1, 2, 575000},
		{1, 1, 580000},
		{2, 2, 85000},
	}
	for _, b := range bids {
		_, err := DB.Exec(
			"INSERT INTO bids (auction_id, user_id, amount, created_at) VALUES (?, ?, ?, ?)",
			b.AuctionID, b.UserID, b.Amount, now,
		)
		if err != nil {
			log.Printf("  Erreur insertion bid: %v", err)
			continue
		}
		DB.Exec("UPDATE auctions SET current_price = ?, winner_id = ? WHERE id = ?", b.Amount, b.UserID, b.AuctionID)
		log.Printf("  Enchère placée: auction #%d, user #%d, %.0f€", b.AuctionID, b.UserID, b.Amount)
	}
}

func seedSalesHistory() {
	sales := []struct {
		PropType string
		Price    float64
		Rooms    int
		Location string
		Surface  float64
	}{
		{"house", 220000, 4, "Lyon", 110},
		{"apartment", 175000, 2, "Paris", 55},
		{"villa", 720000, 5, "Nice", 200},
		{"studio", 82000, 1, "Toulouse", 22},
	}
	for _, s := range sales {
		_, err := DB.Exec(
			"INSERT INTO sales_history (property_type, price, rooms, location, surface) VALUES (?, ?, ?, ?, ?)",
			s.PropType, s.Price, s.Rooms, s.Location, s.Surface,
		)
		if err != nil {
			log.Printf("  Erreur insertion vente: %v", err)
			continue
		}
		log.Printf("  Vente enregistrée: %s %.0f€", s.PropType, s.Price)
	}
}

func seedUserRelations() {
	_, err := DB.Exec("INSERT INTO user_liked_properties (user_id, property_id) VALUES (1, 2), (1, 3)")
	if err != nil {
		log.Printf("  Erreur liked properties: %v", err)
	} else {
		log.Println("  toto aime les propriétés #2, #3")
	}
	_, err = DB.Exec("INSERT INTO user_liked_properties (user_id, property_id) VALUES (2, 1), (2, 5), (2, 7)")
	if err != nil {
		log.Printf("  Erreur liked properties: %v", err)
	} else {
		log.Println("  Alexandre aime les propriétés #1, #5, #7")
	}
	_, err = DB.Exec("INSERT INTO user_bought_properties (user_id, property_id) VALUES (1, 4)")
	if err != nil {
		log.Printf("  Erreur bought properties: %v", err)
	} else {
		log.Println("  toto a acheté la propriété #4")
	}
}
