package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

		var artwork models.Searches
		if err := db.Table("searches").Where("searches.\"ID\" = ?", id).Find(&artwork).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message":       err.Error(),
				"request_param": c.Param("id"),
			})
			log.Print(err.Error())

			return
		}

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
			log.Print(err)
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})

			return
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
			log.Print("[Error] Unable to parse form data")

			return
		}

		pwd := reqData.Password
		password, pwdErr := bcrypt.GenerateFromPassword([]byte(pwd), 14)
		if pwdErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			log.Print("Password could not be encrypted")

			return
		}

		user := models.Users{
			Username: reqData.Username,
			Email:    reqData.Email,
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
			log.Print("[Error] Unable to parse form data")

			return
		}

		var user models.Users
		un, pwd := reqData.Username, reqData.Password
		db.Find(&user, "username = ?", un)
		if user.Username == "" {
			c.Writer.Header().Add("error", "User could not be found through username")
			c.JSON(http.StatusNotFound, gin.H{
				"message": "user could not be found through username",
			})
			log.Print("user could not be found through username")

			return
		}

		pwdErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
		if pwdErr != nil {
			c.Writer.Header().Add("error", "User could not be found through password")
			c.JSON(http.StatusBadRequest, gin.H{
				"password error": pwdErr.Error(),
			})
			log.Print(pwdErr)

			return
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
			log.Print(err)

			return
		}

		c.SetCookie("jwt", token, 60*60*24, "/", "localhost", false, true)

		if user.ID != 0 && pwdErr == nil {
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unable to login",
			})

			return
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
			log.Print("cookie could not be found for user")

			return
		}

		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, keyFunc)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": fmt.Errorf("unauthenticated user: %+v", err),
			})
			log.Printf("unauthenticated user: %+v", err)

			return
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

func GetUsernames(db *gorm.DB) gin.HandlerFunc {
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

		err := reqData.ProcessReq(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			log.Print(err)

			return
		}

		iID, aErr := strconv.Atoi(reqData.ItemID)
		if aErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":         aErr.Error(),
				"reData":          reqData.ToString(),
				"ArtworkID":       reqData.ItemID,
				"post-conversion": iID,
			})
			log.Print(aErr)

			return
		}

		var artworkLike models.ArtworkLikes
		exists, err := likeExists(db, &artworkLike, reqData.UserID, iID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			log.Print(err)

			return
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
		err := reqData.ProcessReq(c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
			})
			log.Print(err.Error())

			return
		}

		iID, aErr := strconv.Atoi(reqData.ItemID)
		if aErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMessage":    aErr.Error(),
				"reData":          reqData.ToString(),
				"ArtworkID":       reqData.ItemID,
				"post-conversion": iID,
			})
			log.Print(aErr)

			return
		}

		var artworkLike models.ArtworkLikes
		exists, err := likeExists(db, &artworkLike, reqData.UserID, iID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorMessage": err.Error(),
			})
			log.Print(err)

			return
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
					"errorMessage": result.Error,
				})
				log.Print(result.Error)

				return

			} else {
				c.JSON(http.StatusCreated, newArtworkLike)
			}
		}
	}
}

func Users(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": errors.Wrap(err, "unable to read request.body").Error(),
			})
			log.Print(err)

			return
		}

		var ID string
		if err := json.Unmarshal(id, &ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message":  errors.Wrap(err, "unable to unmarshal data").Error(),
				"data":     id,
				"ID":       ID,
				"req body": fmt.Sprintf("%v", c.Request.Body),
			})
			log.Print(err)

			return
		}

		var user models.Users
		db.Where("id = ?", ID).Find(&user)
		c.JSON(http.StatusOK, user)
	}
}

func LikedArtworkHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageInt, exist := c.Get("pageInt")
		if !exist {
			msg, _ := c.Get("pageError")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": msg.(error).Error(),
			})
			log.Print(msg)

			return
		}

		userID, exist := c.Get("userID")
		if !exist {
			msg, _ := c.Get("userError")
			c.JSON(http.StatusBadRequest, gin.H{
				"message": msg.(error).Error(),
			})
			log.Print(msg)

			return
		}

		var count int64
		db.Table("artwork_likes").Where("artwork_likes.like = true and user_id = ?", userID.(string)).Offset(pageInt.(int)).Count(&count)

		var likedArtwork []models.Searches
		db.Table("searches").Select(
			"\"searches\".*").Joins(
			"left join artwork_likes as al on al.artwork_id = \"searches\".\"ID\"").Where(
			"al.user_id = ?", userID.(string)).Offset(pageInt.(int)).Scan(&likedArtwork)

		if len(likedArtwork) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message":      "no liked artwork",
				"artworkLikes": 0,
			})
			log.Print("no liked artwork")

			return
		}

		likedList := models.LikedList{
			LikedArtwork: likedArtwork,
			NextPage:     pageInt.(int),
			Count:        count,
		}

		likedList.AddNextPage(10)
		c.JSON(http.StatusOK, likedList)
	}
}
