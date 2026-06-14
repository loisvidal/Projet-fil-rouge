package controllers

import (
	"RedProject/database"
	"RedProject/models"
	"encoding/json"
	"net/http"
	"strconv"
)

var ProfilConnect models.User

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	u, err := database.GetUserByIDWithRelations(StructHome.Profil.IdUser)
	if err != nil {
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
		return
	}

	home, _ := ReloadHome()
	data := struct {
		Profil     models.User
		UserDetail models.User
		Properties []models.Property
	}{
		Profil:     home.Profil,
		UserDetail: u,
		Properties: home.ListProperty,
	}

	temp.ExecuteTemplate(w, "profil", data)
}

func ViewOtherProfil(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
		return
	}

	u, err := database.GetUserByID(id)
	if err != nil {
		http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
		return
	}

	home, _ := ReloadHome()
	data := struct {
		Profil     models.User
		UserDetail models.User
	}{
		Profil:     home.Profil,
		UserDetail: u,
	}

	temp.ExecuteTemplate(w, "profil", data)
}

func MailToOtherProfil(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	StructHome.Profil = models.User{}
	http.Redirect(w, r, "/red_project/home", http.StatusSeeOther)
}

// ============= Admin =============

type AdminPageData struct {
	Profil     models.User
	Users      []models.User
	TotalUsers int
	EditUser   *models.User
}

func AdminUsersHandler(w http.ResponseWriter, r *http.Request) {
	if !isConnectedAdmin() {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Erreur chargement utilisateurs", http.StatusInternalServerError)
		return
	}

	data := AdminPageData{
		Profil:     StructHome.Profil,
		Users:      users,
		TotalUsers: len(users),
	}

	temp.ExecuteTemplate(w, "admin_users", data)
}

func AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !isConnectedAdmin() {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	if id == StructHome.Profil.IdUser {
		http.Error(w, "Vous ne pouvez pas vous supprimer vous-même", http.StatusBadRequest)
		return
	}

	database.DeleteUser(id)
	http.Redirect(w, r, "/red_project/admin/users", http.StatusSeeOther)
}

func AdminToggleAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || !isConnectedAdmin() {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	database.ToggleAdmin(id)
	http.Redirect(w, r, "/red_project/admin/users", http.StatusSeeOther)
}

func AdminEditUser(w http.ResponseWriter, r *http.Request) {
	if !isConnectedAdmin() {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		u := models.User{
			IdUser:   id,
			NameUser: r.FormValue("name"),
			Mail:     r.FormValue("email"),
		}

		password := r.FormValue("password")
		if password != "" {
			hash, err := database.HashPassword(password)
			if err == nil {
				database.UpdatePassword(id, hash)
			}
		}

		database.UpdateUser(u)
		http.Redirect(w, r, "/red_project/admin/users", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	u, err := database.GetUserByID(id)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
		return
	}

	users, _ := database.GetAllUsers()

	data := AdminPageData{
		Profil:     StructHome.Profil,
		Users:      users,
		TotalUsers: len(users),
		EditUser:   &u,
	}
	temp.ExecuteTemplate(w, "admin_users", data)
}

func AdminUsersJSON(w http.ResponseWriter, r *http.Request) {
	if !isConnectedAdmin() {
		http.Error(w, "Accès refusé", http.StatusForbidden)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		http.Error(w, "Erreur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func isConnectedAdmin() bool {
	return StructHome.Profil.IsConnect && StructHome.Profil.IsAdmin
}
