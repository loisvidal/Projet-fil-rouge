package controllers

import (
	models "RedProject/models"
	"fmt"
	"net/http"
	"text/template"
)

var temp *template.Template
var StructHome models.Home

func Init() {
	var err error
	temp, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Erreur lors du chargement des templates :", err)
	}
}

// home

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	home, err := ReloadHome()
	if err != nil {
		fmt.Println(models.Red, "Property loading error : ", err, models.Reset)
		http.Redirect(w, r, "/RedProject/error", http.StatusSeeOther)
		return
	}
	fmt.Println(StructHome.Profil)
	err = temp.ExecuteTemplate(w, "home", home)
	if err != nil {
		fmt.Println(models.Red, "Template error :", err, models.Reset)
	}
}

func FilterHome(w http.ResponseWriter, r *http.Request) {

}

func ReloadHome() (*models.Home, error) {
	var err error

	if !StructHome.Profil.IsConnect {
		StructHome.Profil = ProfilConnect
	}
	StructHome.ListProperty, err = LoadProperties("Data/dataProperty.json", false)
	if err != nil {
		return nil, err
	}
	return &StructHome, nil
}

// header

func Search(w http.ResponseWriter, r *http.Request) {

}

// Error

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "error", nil)
}
