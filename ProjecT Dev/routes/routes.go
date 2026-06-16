package routes

import (
	"net/http"
	controllers "RedProject/controllers"
)

func InitRoutes() {
	http.HandleFunc("/", controllers.HomeRedirect)
	http.HandleFunc("/login", controllers.LoginPageHandler)
	http.HandleFunc("/register", controllers.RegisterPageHandler)

	http.HandleFunc("/red_project/home", controllers.HomeHandler)
	http.HandleFunc("/red_project/home/filter", controllers.FilterHome)

	http.HandleFunc("/red_project/search", controllers.SearchHandler)

	http.HandleFunc("/red_project/register", controllers.Register)
	http.HandleFunc("/red_project/login", controllers.Login)
	http.HandleFunc("/red_project/confirm", controllers.ConfirmEmailHandler)
	http.HandleFunc("/red_project/logout", controllers.LogoutHandler)

	http.HandleFunc("/red_project/logement/", controllers.LogementHandler)
	http.HandleFunc("/red_project/logement/buy", controllers.BuyLogement)
	http.HandleFunc("/red_project/logement/post", controllers.PostLogement)
	http.HandleFunc("/red_project/logement/note", controllers.NoteLogement)

	http.HandleFunc("/red_project/Profil", controllers.ProfilHandler)
	http.HandleFunc("/red_project/Profil/delete", controllers.DeleteOwnAccount)
	http.HandleFunc("/red_project/Profil/Consult/", controllers.ViewOtherProfil)
	http.HandleFunc("/red_project/Profil/mail", controllers.MailToOtherProfil)

	http.HandleFunc("/red_project/auctions", controllers.AuctionHandler)
	http.HandleFunc("/red_project/auction/", controllers.AuctionDetailHandler)
	http.HandleFunc("/red_project/auction/bid", controllers.PlaceBidHandler)
	http.HandleFunc("/red_project/auction/create", controllers.CreateAuctionHandler)

	http.HandleFunc("/red_project/analytics", controllers.AnalyticsHandler)
	http.HandleFunc("/red_project/analytics/predict", controllers.PredictPriceHandler)
	http.HandleFunc("/red_project/analytics/report", controllers.SalesReportHandler)
	http.HandleFunc("/red_project/analytics/geocode", controllers.GeocodeHandler)

	http.HandleFunc("/red_project/admin/users", controllers.AdminUsersHandler)
	http.HandleFunc("/red_project/admin/users/delete", controllers.AdminDeleteUser)
	http.HandleFunc("/red_project/admin/users/toggle-admin", controllers.AdminToggleAdmin)
	http.HandleFunc("/red_project/admin/users/edit", controllers.AdminEditUser)
	http.HandleFunc("/red_project/admin/api/users", controllers.AdminUsersJSON)

	http.HandleFunc("/red_project/error", controllers.ErrorHandler)

	http.HandleFunc("/red_project/legal/", controllers.LegalHandler)
	http.HandleFunc("/red_project/contact", controllers.ContactHandler)
}
