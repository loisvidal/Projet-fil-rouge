package controllers

import (
	"RedProject/database"
	models "RedProject/models"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func LogementHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/red_project/logement/")
	if idStr == "" || idStr == "/" {
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/red_project/error", http.StatusSeeOther)
		return
	}

	prop, err := database.GetPropertyByID(id)
	if err != nil {
		http.Redirect(w, r, "/red_project/error", http.StatusSeeOther)
		return
	}

	auction, _ := database.GetActiveAuctions()
	var propAuction *models.Auction
	for _, a := range auction {
		if a.PropertyID == id {
			propAuction = &a
			break
		}
	}

	home, _ := ReloadHome()
	data := struct {
		Profil   models.User
		Property models.Property
		Auction  *models.Auction
	}{
		Profil:   home.Profil,
		Property: prop,
		Auction:  propAuction,
	}

	temp.ExecuteTemplate(w, "logement", data)
}

func BuyLogement(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	err = database.BuyProperty(StructHome.Profil.IdUser, id)
	if err != nil {
		http.Error(w, "Erreur lors de l'achat", http.StatusInternalServerError)
		return
	}

	prop, _ := database.GetPropertyByID(id)
	database.RecordSale(prop, prop.PriceProperty)

	if StructHome.Profil.Mail != "" {
		subject := "Confirmation d'achat YPlaza"
		body := fmt.Sprintf(`Bonjour %s,

Félicitations ! Vous avez acheté le bien suivant :

%s
Prix : %.0f €
Localisation : %s

Merci de votre confiance.
L'équipe YPlaza`, StructHome.Profil.NameUser, prop.NameProperty, prop.PriceProperty, prop.Location)

		from := "alexandre.petitfrere@ynov.com"
		addr := os.Getenv("SMTP_ADDR")
		host := os.Getenv("SMTP_HOST")
		smtpUser := os.Getenv("SMTP_USER")
		smtpPass := os.Getenv("SMTP_PASS")

		if addr != "" && host != "" {
			auth := smtp.PlainAuth("", smtpUser, smtpPass, host)
			msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", from, StructHome.Profil.Mail, subject, body)
			if err := smtp.SendMail(addr, auth, from, []string{StructHome.Profil.Mail}, []byte(msg)); err != nil {
				fmt.Printf("Erreur envoi confirmation achat à %s: %v\n", StructHome.Profil.Mail, err)
			} else {
				fmt.Printf("Confirmation d'achat envoyée à %s\n", StructHome.Profil.Mail)
			}
		} else {
			fmt.Printf("[ACHAT] De: %s\nPour: %s\n%s\n%s\n", from, StructHome.Profil.Mail, subject, body)
		}
	}

	http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
}

func PostLogement(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}

	if r.Method == http.MethodGet {
		home, _ := ReloadHome()
		temp.ExecuteTemplate(w, "post_logement", home)
		return
	}

	if r.Method != http.MethodPost {
		return
	}

	r.ParseMultipartForm(10 << 20)

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	rooms, _ := strconv.Atoi(r.FormValue("rooms"))
	surface, _ := strconv.ParseFloat(r.FormValue("surface"), 64)

	propType := r.FormValue("property_type")

	imgProperty := database.DefaultImagesForType(propType)

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			ext = ".jpg"
		}
		filename := fmt.Sprintf("upload_%d%s", time.Now().UnixNano(), ext)

		imgDir := assetsImgDir()
		os.MkdirAll(imgDir, os.ModePerm)
		dst, err := os.Create(filepath.Join(imgDir, filename))
		if err == nil {
			defer dst.Close()
			io.Copy(dst, file)
			imgProperty = []string{filename}
		}
	}

	prop := models.Property{
		NameProperty:   r.FormValue("name"),
		DescProprety:   r.FormValue("description"),
		PriceProperty:  price,
		IsSellProperty: true,
		ImgProperty:    imgProperty,
		OwnerID:        StructHome.Profil.IdUser,
		Type:           propType,
		Rooms:          rooms,
		Location:       r.FormValue("location"),
		Surface:        surface,
	}

	_, err = database.CreateProperty(prop)
	if err != nil {
		http.Error(w, "Erreur lors de la création", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
}

func NoteLogement(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
		return
	}

	query := r.FormValue("searchHeader")
	results, err := database.SearchProperties(query)
	if err != nil {
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
		return
	}

	home, _ := ReloadHome()
	data := struct {
		Profil       models.User
		ListProperty []models.Property
		SearchQuery  string
	}{
		Profil:       home.Profil,
		ListProperty: results,
		SearchQuery:  query,
	}

	temp.ExecuteTemplate(w, "home", data)
}

func FilterHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	propType := r.FormValue("type")
	minPrice := r.FormValue("min_price")
	maxPrice := r.FormValue("max_price")
	minSurface := r.FormValue("min_surface")
	maxSurface := r.FormValue("max_surface")
	location := r.FormValue("location")

	allProps, err := database.GetAllProperties()
	if err != nil {
		return
	}

	var filtered []models.Property
	for _, p := range allProps {
		if propType != "" && p.Type != propType {
			continue
		}
		if minPrice != "" {
			min, _ := strconv.ParseFloat(minPrice, 64)
			if p.PriceProperty < min {
				continue
			}
		}
		if maxPrice != "" {
			max, _ := strconv.ParseFloat(maxPrice, 64)
			if p.PriceProperty > max {
				continue
			}
		}
		if minSurface != "" {
			min, _ := strconv.ParseFloat(minSurface, 64)
			if p.Surface < min {
				continue
			}
		}
		if maxSurface != "" {
			max, _ := strconv.ParseFloat(maxSurface, 64)
			if p.Surface > max {
				continue
			}
		}
		if location != "" && !strings.Contains(strings.ToLower(p.Location), strings.ToLower(location)) {
			continue
		}
		filtered = append(filtered, p)
	}

	if filtered == nil {
		filtered = []models.Property{}
	}

	home, _ := ReloadHome()
	data := struct {
		Profil       models.User
		ListProperty []models.Property
	}{
		Profil:       home.Profil,
		ListProperty: filtered,
	}

	temp.ExecuteTemplate(w, "home", data)
}

// ============= Auction Handlers =============

func AuctionHandler(w http.ResponseWriter, r *http.Request) {
	auctions, err := database.GetActiveAuctions()
	if err != nil {
		http.Error(w, "Erreur chargement enchères", http.StatusInternalServerError)
		return
	}

	home, _ := ReloadHome()
	now := time.Now()

	type AuctionItem struct {
		models.Auction
		RemainingSeconds int
		ImgProperty      []string
		TypeLabel        string
		Surface          float64
		Rooms            int
		Location         string
	}

	var auctionsWithTime []AuctionItem

	for _, a := range auctions {
		remaining := int(a.EndTime.Sub(now).Seconds())
		if remaining < 0 {
			remaining = 0
		}
		prop, _ := database.GetPropertyByID(a.PropertyID)
		item := AuctionItem{a, remaining, prop.ImgProperty, prop.TypeLabel(), prop.Surface, prop.Rooms, prop.Location}
		auctionsWithTime = append(auctionsWithTime, item)
	}

	data := struct {
		Profil   models.User
		Auctions []AuctionItem
	}{
		Profil:   home.Profil,
		Auctions: auctionsWithTime,
	}

	temp.ExecuteTemplate(w, "auctions", data)
}

func AuctionDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/red_project/auction/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/red_project/auctions", http.StatusSeeOther)
		return
	}

	auction, err := database.GetAuctionByID(id)
	if err != nil {
		http.Redirect(w, r, "/red_project/auctions", http.StatusSeeOther)
		return
	}

	bids, _ := database.GetAuctionBids(id)
	prop, _ := database.GetPropertyByID(auction.PropertyID)
	home, _ := ReloadHome()

	now := time.Now()
	remaining := int(auction.EndTime.Sub(now).Seconds())
	if remaining < 0 {
		remaining = 0
	}

	data := struct {
		Profil           models.User
		Auction          models.Auction
		Property         models.Property
		Bids             []models.Bid
		RemainingSeconds int
	}{
		Profil:           home.Profil,
		Auction:          auction,
		Property:         prop,
		Bids:             bids,
		RemainingSeconds: remaining,
	}

	temp.ExecuteTemplate(w, "auction_detail", data)
}

func PlaceBidHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		return
	}

	auctionID, _ := strconv.Atoi(r.FormValue("auction_id"))
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)

	if auctionID == 0 || amount <= 0 {
		http.Error(w, "Montant invalide", http.StatusBadRequest)
		return
	}

	auction, err := database.GetAuctionByID(auctionID)
	if err != nil || !auction.IsActive || time.Now().After(auction.EndTime) {
		http.Error(w, "Enchère invalide ou terminée", http.StatusBadRequest)
		return
	}

	minBid := auction.CurrentPrice + auction.MinBidStep
	if amount < minBid {
		http.Error(w, fmt.Sprintf("L'enchère minimale est de %.2f €", minBid), http.StatusBadRequest)
		return
	}

	if err := database.PlaceBid(auctionID, StructHome.Profil.IdUser, amount); err != nil {
		http.Error(w, "Erreur lors de l'enchère", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/red_project/auction/%d", auctionID), http.StatusSeeOther)
}

func CreateAuctionHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}

	if r.Method == http.MethodGet {
		home, _ := ReloadHome()
		myProps, _ := database.GetUserUnsoldProperties(StructHome.Profil.IdUser)
		if myProps == nil {
			myProps = []models.Property{}
		}
		data := struct {
			Profil          models.User
			UserProperties  []models.Property
		}{
			Profil:         home.Profil,
			UserProperties: myProps,
		}
		temp.ExecuteTemplate(w, "create_auction", data)
		return
	}

	if r.Method != http.MethodPost {
		return
	}

	propID, _ := strconv.Atoi(r.FormValue("property_id"))
	startPrice, _ := strconv.ParseFloat(r.FormValue("start_price"), 64)
	minStep, _ := strconv.ParseFloat(r.FormValue("min_bid_step"), 64)
	durationHours, _ := strconv.Atoi(r.FormValue("duration"))

	if propID == 0 || startPrice <= 0 {
		http.Error(w, "Paramètres invalides", http.StatusBadRequest)
		return
	}
	if durationHours <= 0 {
		durationHours = 24
	}

	endTime := time.Now().Add(time.Duration(durationHours) * time.Hour)

	_, err := database.CreateAuction(propID, StructHome.Profil.IdUser, startPrice, minStep, endTime)
	if err != nil {
		http.Error(w, "Erreur création enchère", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/red_project/auctions", http.StatusSeeOther)
}

func assetsImgDir() string {
	cwd, _ := os.Getwd()
	dirs := []string{
		filepath.Join(cwd, "assets", "img"),
	}
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Join(filepath.Dir(exe), "assets", "img"))
	}
	for _, d := range dirs {
		if info, err := os.Stat(d); err == nil && info.IsDir() {
			return d
		}
	}
	return dirs[0]
}
