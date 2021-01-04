package routes

import (
	"os"

	"bitbucket.org/indoquran-api/src/config"
	handle_quran "bitbucket.org/indoquran-api/src/handlers/quran"
	handle_user "bitbucket.org/indoquran-api/src/handlers/user"
	"bitbucket.org/indoquran-api/src/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// StartGin : start the server
func StartGin() {
	r := gin.Default()
	r.Use(cors.Default())
	port := os.Getenv("PORT")
	if port == "80" || port == "" {
		port = config.Config.Server.Port
	}

	// just for testing
	r.GET("/test", handle_quran.TestGet)

	api := r.Group("/api")
	{
		api.Use(middlewares.VisitorLogger)

		usr := api.Group("/user")
		{
			usr.POST("/login", handle_user.LoginUser)
			usr.GET("/:id", handle_user.GetUser)
		}

		// api.GET("/users", handle_user.GetAllUser)
		// api.POST("/users", handle_user.CreateUser)
		// api.GET("/users/:id", handle_user.GetUser)
		// api.PUT("/users/:id", handle_user.UpdateUser)
		// api.DELETE("/users/:id", handle_user.DeleteUser)

		quran := api.Group("/quran")
		{
			quran.GET("/surat", handle_quran.GetSurats)
			quran.GET("/surat/:id", handle_quran.GetSurats)
			quran.GET("/ayat", handle_quran.GetSearchAyats)
			quran.GET("/ayat/:surat_ayat", handle_quran.GetDetailAyat)
			quran.GET("/search/:searchText", handle_quran.GetSearchAyats)
			quran.GET("/kata-bijak", handle_quran.GetKataBijak)

			imp := quran.Group("/import")
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
		}
	}

	r.Run(":" + port)
}
