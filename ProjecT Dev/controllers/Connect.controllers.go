package controllers

import (
	"RedProject/models"
	"fmt"
	"net/http"
)

func sendVerificationEmail(email string) {
	fmt.Printf("Email de vérification envoyé à : %s\n", email)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/red_project/error", http.StatusMovedPermanently)
		return
	}

	name := r.FormValue("userConnect")
	password := r.FormValue("userPassword")
	email := r.FormValue("userMail")

	newUser := models.User{
		NameUser: name,
		Password: password,
		Mail:     email,
	}

	WriteUserConnect(newUser)

	if email != "" {
		sendVerificationEmail(email)
	}
	StructHome.Profil = newUser
	http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	name := r.FormValue("userConnect")
	password := r.FormValue("userPassword")

	if CheckUserConnect(name, password) {
		StructHome.Profil = CheckUser
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
	} else {
		http.Error(w, "Identifiants incorrects", http.StatusUnauthorized)
	}
}
