package main

import (
	han "AT-BE/handlers"
	"AT-BE/utils"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// cors middleware is set to all origins for testing
	c := cors.DefaultConfig()
	c.AllowAllOrigins = true
	router.Use(cors.New(c))

	db, err := utils.SetupConfiguration(false)
	if err != nil {
		e := fmt.Errorf("failed to connect to db %v", err)
		panic(e)
	}

	// fmt.Println("--migrating Users, ArtworkLikes, Curations, CurationLikes--")
	// db.AutoMigrate(&models.Users{}, &models.ArtworkLikes{}, &models.Curations{}, &models.CurationLikes{}, &models.CurationArtwork{})

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
	router.POST("/logout", han.Logout(db))

	router.POST("like", han.ArtworkLike(db))
	router.POST("likes", han.CheckArtworkLikes(db))

	router.Run(os.Getenv("host"))
}
