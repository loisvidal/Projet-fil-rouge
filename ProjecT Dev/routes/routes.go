package routes

import (
	"net/http"
	controllers "RedProject/controllers"
)

func InitRoutes() {
	// home
	http.HandleFunc("/RedProject/home", controllers.HomeHandler) //................... Get
	http.HandleFunc("/RedProject/home/filter", controllers.FilterHome) //............. Post

	// header
	http.HandleFunc("/RedProject/search", controllers.Search) //...................... Post
	
	// connect 
	http.HandleFunc("/RedProject/connect", controllers.ConnectHandler) //............. Get
	http.HandleFunc("/RedProject/register", controllers.Register) //.................. Post
	http.HandleFunc("/RedProject/login", controllers.Login) //........................ Post

	// logement
	http.HandleFunc("/RedProject/logement/", controllers.LogementHandler) //.......... Get
	http.HandleFunc("/RedProject/logement/buy", controllers.BuyLogement) //........... Post
	http.HandleFunc("/RedProject/logement/post", controllers.PostLogement) //......... Post
	http.HandleFunc("/RedProject/logement/note", controllers.NoteLogement) //......... Post

	// Profil
	http.HandleFunc("/RedProject/Profil", controllers.ProfilHandler) //............... Get
	http.HandleFunc("/RedProject/Profil/Consult/", controllers.ViewOtherProfil) //.... Get
	http.HandleFunc("/RedProject/Profil/mail", controllers.MailToOtherProfil) //...... Post

}