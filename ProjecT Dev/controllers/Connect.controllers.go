package controllers

import (
	"RedProject/models"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

func sendConfirmationEmail(email, token string) {
	link := fmt.Sprintf("http://localhost:8080/red_project/confirm?token=%s", token)
	subject := "Confirmation de votre inscription YPlaza"
	body := fmt.Sprintf(`Bonjour,

Merci de vous être inscrit sur YPlaza.

Veuillez confirmer votre adresse email en cliquant sur le lien ci-dessous (valable 2 minutes) :

%s

Ce lien expirera dans 2 minutes.

Cordialement,
L'équipe YPlaza`, link)

	from := os.Getenv("SMTP_FROM")
	addr := os.Getenv("SMTP_ADDR")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")

	if from != "" && addr != "" && host != "" {
		auth := smtp.PlainAuth("", smtpUser, smtpPass, host)
		msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", from, email, subject, body)
		err := smtp.SendMail(addr, auth, from, []string{email}, []byte(msg))
		if err != nil {
			fmt.Printf("Erreur envoi email à %s: %v\n", email, err)
			fmt.Printf("[CONFIRMATION] %s\n%s\n", email, link)
			return
		}
		fmt.Printf("Email de confirmation envoyé à %s\n", email)
		return
	}

	fmt.Printf("[CONFIRMATION] %s\n%s\n", email, link)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/red_project/error", http.StatusMovedPermanently)
		return
	}

	name := strings.TrimSpace(r.FormValue("userConnect"))
	password := r.FormValue("userPassword")
	email := strings.TrimSpace(r.FormValue("userMail"))

	if name == "" || password == "" || email == "" {
		http.Error(w, "Tous les champs sont obligatoires", http.StatusBadRequest)
		return
	}

	if len(password) < 6 {
		http.Error(w, "Mot de passe trop court (min 6 caractères)", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		NameUser: name,
		Password: password,
		Mail:     email,
	}

	token, err := WriteUserConnect(newUser)
	if err != nil {
		http.Error(w, "Erreur lors de l'inscription (nom ou email déjà pris)", http.StatusConflict)
		return
	}

	if email != "" {
		sendConfirmationEmail(email, token)
	}

	newUser.IsConnect = true
	newUser.IsConfirmed = false
	StructHome.Profil = newUser
	http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	identifier := strings.TrimSpace(r.FormValue("userConnect"))
	password := r.FormValue("userPassword")

	if identifier == "" || password == "" {
		http.Error(w, "Identifiants requis", http.StatusBadRequest)
		return
	}

	if CheckUserConnect(identifier, password) {
		StructHome.Profil = CheckUser
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
	} else {
		http.Error(w, "Identifiants incorrects ou compte verrouillé", http.StatusUnauthorized)
	}
}
