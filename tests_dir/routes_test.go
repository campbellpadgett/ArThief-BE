package tests_dir

import (
	"AT-BE/env_var"
	"AT-BE/handlers"
	"AT-BE/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// type Era struct {
// 	ID            int
// 	Era_Name      string
// 	Last_Modified time.Time
// }

func setupGetRouter(handler gin.HandlerFunc, route string) *gin.Engine {
	r := gin.New()
	r.GET(route, handler)
	r.Use(gin.Recovery())

	return r
}

// func TestSignUp(t *testing.T) {
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect to db")
// 	}

// 	route := "/sign-up"
// 	handler := handlers.RegisterUser(db)
// 	router := setupGetRouter(handler, route)
// 	writer := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, route, strings.NewReader("username=cpadgett&email=cpadgett@gmail.com&password=my_Pass11"))
// 	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	router.ServeHTTP(writer, req)

// 	assert.Equal(t, 201, writer.Code)

// 	data, err := ioutil.ReadAll(writer.Body)
// 	if err != nil {
// 		t.Errorf("[Error] Unable to read writer.Body: %s", err)
// 	}

// 	var user models.Users
// 	if err := json.Unmarshal(data, &user); err != nil {
// 		t.Errorf("[ERROR] Unable to unmarshal data to user: %s", err)
// 	}

// 	assert.True(t, true, user.Username == "cpadgett")
// 	fmt.Printf("%v+", user)
// }

func TestUsernames(t *testing.T) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	route := "/usernames"
	handler := handlers.Usernames(db)
	router := setupGetRouter(handler, route)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, route, nil)

	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var usernames map[string]bool
	if err := json.Unmarshal(data, &usernames); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to usernames: %s", err)
	}

	fmt.Println(usernames)
	assert.True(t, true, len(usernames) >= 1)
}

func TestSearch(t *testing.T) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	route := "/search/:term"
	handler := handlers.Search(db)
	router := setupGetRouter(handler, route)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/search/jacob", nil)

	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var searches []models.Searches
	if err := json.Unmarshal(data, &searches); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to searches: %s", err)
	}

	assert.True(t, true, len(searches) > 1)
	assert.Nil(t, nil, searches[0].IMG_S)
	assert.True(t, true, searches[0].ID == "22")

	// fmt.Printf("%v+", searches)
}

func TestGetEra(t *testing.T) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	route := "/era/:id"
	handler := handlers.GetEra(db)
	router := setupGetRouter(handler, route)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/era/5", nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var era models.Era
	if err := json.Unmarshal(data, &era); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to era: %s", err)
	}

	assert.Equal(t, "1817", era.Era_Name)
}

func TestGetSource(t *testing.T) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	route := "/source/:id"
	handler := handlers.GetSource(db)
	router := setupGetRouter(handler, route)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/source/5", nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var source models.Source
	if err := json.Unmarshal(data, &source); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to source: %s", err)
	}

	assert.Equal(t, "Wilson L. Mead Fund", source.Source_Name)
	assert.Equal(t, "CHI", source.Abbreviation)
}

func TestGetArtist(t *testing.T) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	route := "/artist/:id"
	handler := handlers.GetArtist(db)
	router := setupGetRouter(handler, route)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/artist/5", nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var artist models.Artist
	if err := json.Unmarshal(data, &artist); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to artist: %s", err)
	}

	assert.Equal(t, "Ishizaki Yushi", artist.Name)
	assert.Equal(t, "1817", artist.Era)
	assert.Equal(t, " Japan", artist.Description)
	assert.Equal(t, "Male", artist.Gender)
}

func TestGetArtwork(t *testing.T) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	route := "/artwork/:id"
	handler := handlers.GetArtwork(db)
	router := setupGetRouter(handler, route)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/artwork/5", nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var artwork models.Artwork
	if err := json.Unmarshal(data, &artwork); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to artwork: %s", err)
	}

	assert.Equal(t, "Portrait of the Cock Blomhoff family", artwork.Title)
	assert.Equal(t, "1817", artwork.Date_of_Release)
	assert.Equal(t, " Japan", artwork.Artist_Bio)
	assert.Equal(t, "https://lh6.ggpht.com/VK4zTxsKnaZnfOgNyJsMYtt-wf1aGV8rdpIlYCQs4azxhuo_Go3VXkAKR9INbOS5v5v2bREOnlQolXrmK6dsznV3VCw=s0", artwork.Image)
}

func TestGetArtworks(t *testing.T) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", env_var.Host, env_var.User, env_var.Password, env_var.DBname, env_var.Port, env_var.SSLMODE, env_var.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to db")
	}

	route := "/artworks/"
	handler := handlers.GetArtworks(db)
	router := setupGetRouter(handler, route)
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/artworks/?limit=10&last_id=1", nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("\u001b[31m[Error] Unable to read writer.Body: %s", err)
	}

	var artworks []models.Artwork
	if err := json.Unmarshal(data, &artworks); err != nil {
		t.Errorf("\u001b[31m[ERROR] Unable to unmarshal data to artworks: %s", err)
	}

	assert.Equal(t, 10, len(artworks))

	a0, a2, a5, a7 := artworks[0], artworks[2], artworks[5], artworks[7]
	assert.Equal(t, a0.Title, "Portrait of Leonardus van der Voort")
	assert.Equal(t, a2.Title, "Ladies company looks at stereoscopic photos")
	assert.Equal(t, a5.Title, "The harbor entrance of Willemstad with the Government Palace")
	assert.Equal(t, a7.Title, "Heemskerck and Barents prepare their second expedition to the North")
}
