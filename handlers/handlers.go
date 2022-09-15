package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"AT-BE/models"
	"AT-BE/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetEra(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var era models.Era
		db.Where("id = ?", id).Find(&era)

		c.JSON(http.StatusOK, era)
	}
}

func GetArtist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var artist models.Artist
		db.Where("id = ?", id).Find(&artist)

		c.JSON(http.StatusOK, artist)
	}
}

func GetArtwork(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var artwork models.Artwork
		db.Where("id = ?", id).Find(&artwork)

		c.JSON(http.StatusOK, artwork)
	}
}

func GetSource(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var source models.Source
		db.Where("id = ?", id).Find(&source)

		c.JSON(http.StatusOK, source)
	}
}

func GetArtworks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			panic(err)
		}
		last_id := c.Query("last_id")

		var artworks []models.Artwork
		db.Where("id > ?", last_id).Limit(limit).Find(&artworks)

		c.JSON(http.StatusOK, artworks)
	}
}

func Search(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		term := c.Param("term")
		var results []models.Searches
		db.Table("searches").Where("to_tsvector(\"searches\".\"Title\" || '' || \"searches\".\"Artist_Name\") @@ plainto_tsquery(?)", term).Find(&results)

		c.JSON(http.StatusOK, results)
	}
}

func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		reqData, err := utils.ParseFormData(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			panic("[Error] Unable to parse form data")
		}

		pwd := reqData["password"]
		password, pwdErr := bcrypt.GenerateFromPassword([]byte(pwd), 14)
		if pwdErr != nil {
			fmt.Printf("[ERROR] %v", pwdErr)
			panic("Password could not be encrypted")
		}

		user := models.Users{
			Username: reqData["username"],
			Email:    reqData["email"],
			Password: password,
		}

		if pwdErr == nil && err == nil {
			db.Create(&user)
			c.JSON(http.StatusCreated, user)
		}

	}
}

func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		reqData, parseErr := utils.ParseFormData(c.Request.Body)
		if parseErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": parseErr.Error(),
			})
			panic("[Error] Unable to parse form data")
		}

		var user models.Users
		un, pwd := reqData["username"], reqData["password"]
		db.Find(&user, "username = ?", un)
		if user.Username == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user could not be found through username",
			})
			panic("user could not be found through username")
		}

		pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
		if pwdErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"password error": pwdErr.Error(),
			})
			panic(pwdErr)
		}

		stringID := strconv.Itoa(int(user.ID))
		day_length := time.Now().Add(time.Hour * 24)
		claim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    stringID,
			ExpiresAt: day_length.Unix(),
		})

		token, err := claim.SignedString([]byte(os.Getenv("secretkey")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err,
			})
			panic(err)
		}

		c.SetCookie("jwt", token, 60*60*24, "/", "localhost", false, true)

		if user.ID != 0 && pwdErr == nil {
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unable to login",
			})
		}
	}
}

func AuthenticateUser(db *gorm.DB) gin.HandlerFunc {

	// used as arg in jwt.ParseWithClaims below
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("secretkey")), nil
	}

	return func(c *gin.Context) {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(http.StatusNotFound, "cookie could not be found for user")
		}

		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, keyFunc)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": fmt.Errorf("unauthenticated user: %+v", err),
			})
		}

		// has no Issuer attribute due to Claims being an interface, need to type cast
		claim := token.Claims.(*jwt.StandardClaims)

		var user models.Users
		db.Find(&user, "id = ?", claim.Issuer)

		c.JSON(http.StatusOK, user)
	}
}

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("jwt", "", -1, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{
			"message": "successsfully logged out",
		})
	}
}

func Usernames(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var users []models.Users

		db.Find(&users)

		m := make(map[string]bool)

		for _, u := range users {
			m[u.Username] = true
		}

		c.JSON(http.StatusOK, m)
	}
}

// Checks to see if there is already an instance in the db of a user's likeing of artwork
func likeExists(db *gorm.DB, al *models.ArtworkLikes, uID int, aID int) (bool, error) {
	result := db.Where(&models.ArtworkLikes{Artwork_ID: aID, User_ID: uID}, "artwork_id", "user_id").First(&al)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		} else {
			return false, result.Error
		}
	}
	return true, nil
}

// Takes the request data from an [id].tsx page and sends any existing ArtworkLike data and a boolean
func CheckArtworkLikes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var reqData models.LikeReqData

		err := reqData.ProcessLikeReq(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			panic(err)
		}

		iID, aErr := strconv.Atoi(reqData.ItemID)
		if aErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":         aErr.Error(),
				"reData":          reqData.ToString(),
				"ArtworkID":       reqData.ItemID,
				"post-conversion": iID,
			})
			panic(aErr)
		}

		var artworkLike models.ArtworkLikes
		exists, err := likeExists(db, &artworkLike, reqData.UserID, iID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			panic(err)
		}

		if exists {
			c.JSON(http.StatusAccepted, gin.H{
				"data":  artworkLike,
				"liked": artworkLike.Like,
				"exist": true,
			})
		} else {
			c.JSON(http.StatusAccepted, gin.H{
				"data":  nil,
				"liked": false,
				"exist": false,
			})
		}

	}
}

// Takes request data, checks if the data exists in the db and either creates a new instance in the db or updates
// the already existing one.
func ArtworkLike(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var reqData models.LikeReqData
		err := reqData.ProcessLikeReq(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}

		iID, aErr := strconv.Atoi(reqData.ItemID)
		if aErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":         aErr.Error(),
				"reData":          reqData.ToString(),
				"ArtworkID":       reqData.ItemID,
				"post-conversion": iID,
			})

			panic(aErr)
		}

		var artworkLike models.ArtworkLikes
		exists, err := likeExists(db, &artworkLike, reqData.UserID, iID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})

			panic(err)
		}

		if exists {
			artworkLike.Like = reqData.LikeStatus
			db.Save(&artworkLike)

			c.JSON(http.StatusOK, gin.H{
				"message": "successfully updated",
				"like":    artworkLike,
			})
		} else {
			newArtworkLike := models.ArtworkLikes{
				Artwork_ID: iID,
				User_ID:    reqData.UserID,
				Like:       reqData.LikeStatus,
			}

			result := db.Create(&newArtworkLike)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": result.Error,
				})

				panic(result.Error)
			} else {
				c.JSON(http.StatusCreated, newArtworkLike)
			}
		}
	}
}
