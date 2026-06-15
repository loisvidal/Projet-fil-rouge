package controllers

import (
	"RedProject/database"
	"RedProject/models"
	"log"
	"net/http"
	"time"
)

var CheckUser models.User

func CheckUserConnect(identifier string, password string) bool {
	if u, ok := hardcodedCheckUserConnect(identifier, password); ok {
		StructHome.Profil = u
		return true
	}

	u, err := database.GetUserByLogin(identifier)
	if err != nil {
		return false
	}

	if u.LockedUntil != nil && u.LockedUntil.After(time.Now()) {
		log.Printf("Compte %s verrouillé jusqu'à %v", identifier, u.LockedUntil)
		return false
	}

	if !database.CheckPassword(u.Password, password) {
		database.IncrementFailedAttempts(u.IdUser)
		return false
	}

	database.ResetFailedAttempts(u.IdUser)
	u.IsConnect = true
	CheckUser = u
	CheckUser.IsConnect = true
	return true
}

func WriteUserConnect(newUser models.User) (string, error) {
	id, token, err := database.CreateUser(newUser)
	if err != nil {
		return "", err
	}
	newUser.IdUser = int(id)
	return token, nil
}

func requireLogin(w http.ResponseWriter, r *http.Request) bool {
	if !StructHome.Profil.IsConnect {
		isAJAX := r.Header.Get("X-Requested-With") == "XMLHttpRequest" || r.Header.Get("Accept") == "application/json"
		if isAJAX {
			http.Error(w, "Vous devez être connecté pour effectuer cette action", http.StatusUnauthorized)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		return false
	}
	if !StructHome.Profil.IsConfirmed {
		isAJAX := r.Header.Get("X-Requested-With") == "XMLHttpRequest" || r.Header.Get("Accept") == "application/json"
		if isAJAX {
			http.Error(w, "Vous devez confirmer votre email pour effectuer cette action", http.StatusForbidden)
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		return false
	}
	return true
}

func ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token manquant", http.StatusBadRequest)
		return
	}

	err := database.ConfirmUser(token)
	if err != nil {
		log.Printf("Erreur confirmation: %v", err)
		http.Error(w, "Lien de confirmation invalide ou expiré", http.StatusGone)
		return
	}

	if StructHome.Profil.IsConnect {
		StructHome.Profil.IsConfirmed = true
	}

	http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
}
