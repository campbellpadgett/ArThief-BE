package main

import (
	han "AT-BE/handlers"
	"AT-BE/models"
	"AT-BE/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func corsMiddleware(origins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()

	db, origins, err := utils.SetupConfiguration(false)
	if err != nil {
		e := fmt.Errorf("failed to connect to db %v", err)
		panic(e)
	}

	// cors middleware is set to all origins for testing, prod is set to FE host
	// c := cors.DefaultConfig()

	log.Printf("Origins: %s", origins)
	router.Use(corsMiddleware(origins))

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{origins},
	// 	AllowMethods:     []string{"PUT", "POST", "GET", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
	// 	AllowAllOrigins:  false,
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	// c.AllowOrigins = []string{origins}
	// router.Use(cors.New(c))

	fmt.Println("--migrating Users, ArtworkLikes, Curations, CurationLikes--")
	db.AutoMigrate(&models.Users{}, &models.ArtworkLikes{}, &models.Curations{}, &models.CurationLikes{}, &models.CurationArtwork{})

	router.GET("/artwork/:id", han.GetArtwork(db))
	router.GET("/artworks/", han.GetArtworks(db))
	router.GET("/artist/:id", han.GetArtist(db))
	router.GET("/era/:id", han.GetEra(db))
	router.GET("/source/:id", han.GetSource(db))

	router.GET("/search/:term", han.Search(db))
	router.GET("/usernames", han.Usernames(db))

	router.POST("/sign-up", han.RegisterUser(db))
	router.POST("/login", han.LoginUser(db))
	router.GET("/user", han.AuthenticateUser(db))
	router.POST("/users", han.Users(db))
	router.POST("/logout", han.Logout(db))

	router.POST("like", han.ArtworkLike(db))
	router.POST("likes", han.CheckArtworkLikes(db))

	d := fmt.Sprint(os.Getenv("HOST") + ":" + os.Getenv("PORT"))
	router.Run(d)
}
