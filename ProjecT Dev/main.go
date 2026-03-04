package main

import (
	controllers "RedProject/controllers"
	routes "RedProject/routes"
	"log"
	"net/http"
)

func main() {

	controllers.Init()
	routes.InitRoutes()
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Println("Serveur lancé sur http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erreur serveur : %v", err)
	}
}
