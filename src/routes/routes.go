package routes

import (
	"os"

	"bitbucket.org/indoquran-api/src/config"
	"bitbucket.org/indoquran-api/src/handlers"
	handle_email "bitbucket.org/indoquran-api/src/handlers/email"
	handle_quran "bitbucket.org/indoquran-api/src/handlers/quran"
	handle_user "bitbucket.org/indoquran-api/src/handlers/user"
	"bitbucket.org/indoquran-api/src/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// StartGin : start the server
func StartGin() {
	r := gin.Default()
	r.Use(cors.Default())

	store := config.SessionWithCookies()
	r.Use(sessions.Sessions(config.Config.Session.Name, store))

	port := os.Getenv("PORT")
	if port == "80" || port == "" {
		port = config.Config.Server.Port
	}

	api := r.Group("/api")
	{
		api.Use(middlewares.VisitorLogger)

		usr := api.Group("/user")
		{
			usr.POST("/login", handle_user.LoginUser)
			usr.POST("/logout", handle_user.LogoutUser)
			usr.GET("/:id", handle_user.GetUser)
		}

		quran := api.Group("/quran")
		{
			quran.GET("/surat", handle_quran.GetSurats)
			quran.GET("/surat/:id", handle_quran.GetSurats)
			quran.GET("/ayat", handle_quran.GetSearchAyats)
			quran.GET("/ayat/:surat_ayat", handle_quran.GetDetailAyat)
			quran.GET("/search/:searchText", handle_quran.GetSearchAyats)
			quran.GET("/kata-bijak", handle_quran.GetKataBijak)
			quran.GET("/topik/:id", handle_quran.GetTopik)
		}

		email := api.Group("/email")
		{
			email.POST("", handle_email.SendEmail)
		}
	}

	imp := r.Group("/import")
	{
		// imp.GET("/surat", handle_quran.ImportSurat)
		// imp.GET("/ayat", handle_quran.ImportAyat)
		// imp.GET("/catatan", handle_quran.ImportCatatan)
		// imp.GET("/tafsir", handle_quran.ImportTafsir)
		// imp.GET("/juz", handle_quran.ImportJuz)
		// imp.GET("/tafsir/move", handle_quran.ImportTafsirMove)
		// imp.GET("/image", handle_quran.ImportImage)
		// imp.GET("/arab-text", handle_quran.ImportArabText)
		// imp.GET("/wise-words", handle_quran.ImportWiseWords)
		imp.GET("/ayat-id", handle_quran.AddAyatID)
	}

	r.Run(":" + port)
}
