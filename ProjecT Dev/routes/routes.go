package routes

import (
	"net/http"
	controllers "RedProject/controllers"
)

func InitRoutes() {
	// home
	http.HandleFunc("/red_project/home", controllers.HomeHandler) //................... Get
	http.HandleFunc("/red_project/home/filter", controllers.FilterHome) //............. Post

	// header
	http.HandleFunc("/red_project/search", controllers.Search) //...................... Post
	
	// connect 
	http.HandleFunc("/red_project/register", controllers.Register) //.................. Post
	http.HandleFunc("/red_project/login", controllers.Login) //........................ Post

	// logement
	http.HandleFunc("/red_project/logement/", controllers.LogementHandler) //.......... Get
	http.HandleFunc("/red_project/logement/buy", controllers.BuyLogement) //........... Post
	http.HandleFunc("/red_project/logement/post", controllers.PostLogement) //......... Post
	http.HandleFunc("/red_project/logement/note", controllers.NoteLogement) //......... Post

	// Profil
	http.HandleFunc("/red_project/Profil", controllers.ProfilHandler) //............... Get
	http.HandleFunc("/red_project/Profil/Consult/", controllers.ViewOtherProfil) //.... Get
	http.HandleFunc("/red_project/Profil/mail", controllers.MailToOtherProfil) //...... Post

	// Error
	http.HandleFunc("/red_project/error", controllers.ErrorHandler) //................. Get

}