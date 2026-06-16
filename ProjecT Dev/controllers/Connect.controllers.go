package controllers

import (
	"RedProject/database"
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

	from := "alexandre.petitfrere@ynov.com"
	addr := os.Getenv("SMTP_ADDR")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")

	if addr != "" && host != "" {
		auth := smtp.PlainAuth("", smtpUser, smtpPass, host)
		msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s", from, email, subject, body)
		err := smtp.SendMail(addr, auth, from, []string{email}, []byte(msg))
		if err != nil {
			fmt.Printf("Erreur envoi email à %s: %v\n", email, err)
			fmt.Printf("[CONFIRMATION] %s\n%s\n", email, link)
			return
		}
		fmt.Printf("Email de confirmation envoyé à %s depuis %s\n", email, from)
		return
	}

	fmt.Printf("[CONFIRMATION] De: %s\nPour: %s\n%s\n", from, email, link)
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/red_project/error", http.StatusMovedPermanently)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	confirm := r.FormValue("confirm_password")

	if username == "" || email == "" || password == "" {
		RegisterPageHandler(w, r)
		return
	}

	if password != confirm {
		RegisterPageHandler(w, r)
		return
	}

	if len(password) < 6 {
		RegisterPageHandler(w, r)
		return
	}

	_, token, err := database.CreateUser(models.User{
		NameUser: username,
		Mail:     email,
		Password: password,
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "UNIQUE") {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	database.ConfirmUser(token)

	u, _ := database.GetUserByLogin(username)
	u.IsConnect = true
	StructHome.Profil = u
	http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	identifier := strings.TrimSpace(r.FormValue("identifier"))
	password := r.FormValue("password")

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
