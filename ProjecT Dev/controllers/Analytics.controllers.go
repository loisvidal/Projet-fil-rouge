package controllers

import (
	"RedProject/database"
	models "RedProject/models"
	"RedProject/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	summary, err := database.GetAnalyticSummary()
	if err != nil {
		http.Error(w, "Erreur chargement analyses", http.StatusInternalServerError)
		return
	}

	sales, _ := database.GetSalesHistory()
	home, _ := ReloadHome()

	prediction := database.PredictPrice("house", 3, 70.0, "")

	data := struct {
		Profil     models.User
		Summary    *models.AnalyticSummary
		Sales      []models.SalesRecord
		Prediction *models.PricePrediction
		Types      []string
		TypeLabels map[string]string
	}{
		Profil:     home.Profil,
		Summary:    summary,
		Sales:      sales,
		Prediction: prediction,
		Types:      models.PropertyTypes,
		TypeLabels: models.PropertyTypeLabels,
	}

	temp.ExecuteTemplate(w, "analytics", data)
}

func PredictPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST requis", http.StatusMethodNotAllowed)
		return
	}

	propType := r.FormValue("property_type")
	rooms, _ := strconv.Atoi(r.FormValue("rooms"))
	surface, _ := strconv.ParseFloat(r.FormValue("surface"), 64)
	location := r.FormValue("location")
	postalCode := r.FormValue("postal_code")

	if propType == "" {
		propType = "house"
	}
	if rooms <= 0 {
		rooms = 3
	}
	if surface <= 0 {
		surface = 70
	}

	if postalCode != "" && location == "" {
		geoLoc, _ := services.GeocodePostalCode(postalCode, "France")
		if geoLoc != nil && geoLoc.Verified {
			location = geoLoc.City
		}
	}

	result := database.PredictPrice(propType, rooms, surface, location)

	if postalCode != "" && result.PostalCode == "" {
		result.PostalCode = postalCode
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GeocodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST requis", http.StatusMethodNotAllowed)
		return
	}

	city := r.FormValue("city")
	postalCode := r.FormValue("postal_code")
	country := r.FormValue("country")

	if country == "" {
		country = "France"
	}

	var result *services.VerifiedLocation
	var err error

	if postalCode != "" {
		result, err = services.GeocodePostalCode(postalCode, country)
	} else {
		result, err = services.GeocodeAddress(city, "", country)
	}

	if err != nil {
		http.Error(w, "Erreur géocodage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func SalesReportHandler(w http.ResponseWriter, r *http.Request) {
	sales, err := database.GetSalesHistory()
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
		return
	}

	summary, _ := database.GetAnalyticSummary()

	w.Header().Set("Content-Type", "application/json")

	report := map[string]interface{}{
		"summary": summary,
		"sales":   sales,
		"exported_at": fmt.Sprintf("%d échantillons analysés", len(sales)),
	}

	json.NewEncoder(w).Encode(report)
}
