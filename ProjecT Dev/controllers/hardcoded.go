package controllers

import (
	"RedProject/database"
	"RedProject/models"
	"fmt"
	"log"
)

var hardcodedUsers = []models.User{
	{
		IdUser:      1,
		NameUser:    "toto",
		Mail:        "natsu.petitfrere27@gmail.com",
		Password:    "$2a$10$IQWvbLQgWtCBCQaEYBrquu2IAV8RekXJ6fj2.d07hgoPYKPGYMDnK",
		IsConfirmed: true,
		IsConnect:   false,
		IsAdmin:     false,
	},
	{
		IdUser:      2,
		NameUser:    "Alexandre",
		Mail:        "alexandre.petitfrere14@gmail.com",
		Password:    "$2a$10$aLPepRm1lJCCZQNQCnZFhuG/nL1qgw6NkVNizP8b3y4SEgkdZnSNC",
		IsConfirmed: true,
		IsConnect:   false,
		IsAdmin:     true,
	},
}

func findHardcodedUser(identifier string) *models.User {
	for i := range hardcodedUsers {
		u := &hardcodedUsers[i]
		if u.NameUser == identifier || u.Mail == identifier {
			return u
		}
	}
	return nil
}

func findHardcodedUserByID(id int) *models.User {
	for i := range hardcodedUsers {
		u := &hardcodedUsers[i]
		if u.IdUser == id {
			return u
		}
	}
	return nil
}

func isHardcodedUser(id int) bool {
	return findHardcodedUserByID(id) != nil
}

// Override database auth: check hardcoded first, skip DB for users
func hardcodedCheckUserConnect(identifier, password string) (models.User, bool) {
	u := findHardcodedUser(identifier)
	if u == nil {
		return models.User{}, false
	}
	if !database.CheckPassword(u.Password, password) {
		log.Printf("Échec connexion %s: mot de passe incorrect", identifier)
		database.IncrementFailedAttempts(u.IdUser)
		return models.User{}, false
	}
	database.ResetFailedAttempts(u.IdUser)
	u.IsConnect = true
	CheckUser = *u
	CheckUser.IsConnect = true
	CheckUser.IsConfirmed = u.IsConfirmed
	CheckUser.IsAdmin = u.IsAdmin
	return CheckUser, true
}

func hardcodedGetUserByID(id int) (models.User, error) {
	u := findHardcodedUserByID(id)
	if u != nil {
		return *u, nil
	}
	return models.User{}, fmt.Errorf("utilisateur %d introuvable", id)
}

func hardcodedIsAdmin(userID int) bool {
	u := findHardcodedUserByID(userID)
	return u != nil && u.IsAdmin
}
