package database

import (
	"encoding/json"
	"log"
)

func DefaultImagesForType(propType string) []string {
	mapping := map[string][]string{
		"house":      {"/assets/img/maison-1.webp"},
		"apartment":  {"/assets/img/appartement-1.webp"},
		"studio":     {"/assets/img/studio-1.jpg"},
		"loft":       {"/assets/img/appartement-1.webp"},
		"villa":      {"/assets/img/villa-1.jpeg"},
		"townhouse":  {"/assets/img/maison-1.webp"},
		"penthouse":  {"/assets/img/appartement-1.webp"},
		"commercial": {"/assets/img/commerce-1.webp"},
		"land":       {"/assets/img/commerce-1.webp"},
	}
	if imgs, ok := mapping[propType]; ok {
		return imgs
	}
	return []string{"appartement-1.webp"}
}

func EnsurePropertyImages() {
	rows, err := DB.Query("SELECT id, property_type FROM properties WHERE images IS NULL OR images = 'null' OR images = '[]'")
	if err != nil {
		log.Printf("Erreur vérification images: %v", err)
		return
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var id int
		var propType string
		if err := rows.Scan(&id, &propType); err != nil {
			continue
		}
		images := DefaultImagesForType(propType)
		imagesJSON, _ := json.Marshal(images)
		DB.Exec("UPDATE properties SET images = ? WHERE id = ?", string(imagesJSON), id)
		count++
		log.Printf("  Image ajoutée propriété #%d (%s): %v", id, propType, images)
	}
	if count > 0 {
		log.Printf("%d propriétés mises à jour avec des images par défaut.", count)
	} else {
		log.Println("Toutes les propriétés ont déjà des images.")
	}
}
