package controllers

import (
	"RedProject/database"
	models "RedProject/models"
	"fmt"
	"html/template"
	"math"
	"net/http"
)

var temp *template.Template
var StructHome models.Home

func Init() {
	funcMap := template.FuncMap{
		"add": func(a, b float64) float64 {
			return a + b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"round": func(a float64) float64 {
			return math.Round(a)
		},
	}

	var err error
	temp, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Erreur lors du chargement des templates :", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	database.CloseExpiredAuctions()

	home, err := ReloadHome()
	if err != nil {
		fmt.Println(models.Red, "Property loading error : ", err, models.Reset)
		http.Redirect(w, r, "/red_project/error", http.StatusSeeOther)
		return
	}
	err = temp.ExecuteTemplate(w, "home", home)
	if err != nil {
		fmt.Println(models.Red, "Template error :", err, models.Reset)
	}
}

func ReloadHome() (*models.Home, error) {
	if !StructHome.Profil.IsConnect {
		StructHome.Profil = ProfilConnect
	}
	properties, err := database.GetAllProperties()
	if err != nil {
		return nil, err
	}
	StructHome.ListProperty = properties
	return &StructHome, nil
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "error", nil)
}
