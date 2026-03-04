package controllers

import (
	"fmt"
	"net/http"
	"text/template"
)

var temp *template.Template

func Init() {
	var err error
	temp, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Erreur lors du chargement des templates :", err)
	}
}

// home

func HomeHandler(w http.ResponseWriter, r *http.Request) {

}

func FilterHome(w http.ResponseWriter, r *http.Request) {

}

// header

func Search(w http.ResponseWriter, r *http.Request) {

}
