package controllers

import (
	models "RedProject/models"
	"encoding/json"
	"net/http"
	"os"
)

func LogementHandler(w http.ResponseWriter, r *http.Request) {

}

func BuyLogement(w http.ResponseWriter, r *http.Request) {

}

func PostLogement(w http.ResponseWriter, r *http.Request) {

}

func NoteLogement(w http.ResponseWriter, r *http.Request) {

}

// ============= Recup Property =============

func LoadProperties(filename string, filter bool) ([]models.Property, error) {
	var properties []models.Property

	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &properties)
	if err != nil {
		return nil, err
	}

	return properties, nil
}
