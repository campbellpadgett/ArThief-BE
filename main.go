package main

import (
	"errors"
	"fmt"
	"os"

	han "AT-BE/handlers"
	"AT-BE/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()

	// cors middleware is set to all origins for testing
	c := cors.DefaultConfig()
	c.AllowAllOrigins = true
	router.Use(cors.New(c))

	var config utils.Config

	// checking if yaml file with config declartions exists, otherwise, use the env varibale provided by railway
	if _, err := os.Stat("config.yaml"); err == nil {
		config.SetUpViper("config", ".", "yaml")
	} else if errors.Is(err, os.ErrNotExist) {
		config.SetUpRailway()
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Server.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.DBname,
		config.Database.DBport,
		config.Database.SSLMODE,
		config.Database.TimeZone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
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

	router.Run("localhost:8080")
}
