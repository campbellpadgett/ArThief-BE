package tests

import (
	"AT-BE/handlers"
	m "AT-BE/middleware"
	"AT-BE/models"
	"AT-BE/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// POST and GET methods currently available
func setupGetRouter(handler gin.HandlerFunc, route string, httpTest string) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)

	switch httpTest {
	case "GET":
		r.GET(route, handler)
	case "POST":
		r.POST(route, handler)
	}
	r.Use(gin.Recovery())

	return r
}

func TestUpdateCurationName(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/curation/new"
	handler := handlers.NewCurationHandler(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	curReq := models.NewCurationReq{
		Name:      "-*-test curation cpadgett-*-",
		UserID:    16,
		ArtworkID: 1015,
	}

	marshalledData, err := json.Marshal(curReq)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	assert.Equal(t, 201, writer.Code)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	type JsonData struct {
		ID      int
		Message string
	}

	var jsonData JsonData
	if err := json.Unmarshal(wb, &jsonData); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	assert.True(t, jsonData.Message == "success")

	route = "/curation/update"
	handler = handlers.UpdateCurationNameHandler(db)
	router = setupGetRouter(handler, route, "POST")
	writer = httptest.NewRecorder()

	upd := models.UpdateCurName{
		ID:   jsonData.ID,
		Name: curReq.Name,
	}

	marshalledData, err = json.Marshal(upd)
	if err != nil {
		t.Error(err)
	}

	req = httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	assert.Equal(t, 202, writer.Code)

	wb, err = ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(wb, &res); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to res: %s", err)
	}

	assert.True(t, res["message"] == "curation name updated")
	assert.True(t, res["new name"] == upd.Name)
}

func TestNewAndDeleteCuration(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/curation/new"
	handler := handlers.NewCurationHandler(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	curReq := models.NewCurationReq{
		Name:      "-*-test curation cpadgett-*-",
		UserID:    16,
		ArtworkID: 1015,
	}

	marshalledData, err := json.Marshal(curReq)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	assert.Equal(t, 201, writer.Code)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(wb, &jsonData); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	assert.True(t, jsonData["message"] == "success")

	route = "/curation/delete"
	handler = handlers.DeleteCurationHandler(db)
	router = setupGetRouter(handler, route, "POST")
	writer = httptest.NewRecorder()

	marshalledData, err = json.Marshal(jsonData["ID"])
	if err != nil {
		t.Error(err)
	}

	req = httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	assert.Equal(t, 202, writer.Code)

	wb, err = ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var msg map[string]string
	if err := json.Unmarshal(wb, &msg); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to msg: %s", err)
	}

	assert.True(t, msg["message"] == "curation deleted")
}

func TestLikedArtworkHandler(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.GET("/likedArtwork", m.Paginate, handlers.LikedArtworkHandler(db))
	writer := httptest.NewRecorder()

	// sampleUser ID
	route := "/likedArtwork?page=0&userID=16"
	req := httptest.NewRequest(http.MethodGet, route, nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 201, writer.Code)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(wb, &jsonData); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	fmt.Printf("| jsonData: %+v", jsonData)

	var page, count float64
	page, count = 10, 2

	assert.True(t, jsonData["page"] == page)
	assert.True(t, jsonData["count"] == count)

}

func TestLogout(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/login"
	handler := handlers.LoginUser(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	loginReq := utils.ParsedUserRequestData{
		Email:    "sample@gmail.com",
		Username: "sampleUser",
		Password: "sampleUser",
	}

	marshalledData, err := json.Marshal(loginReq)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	route = "/logout"
	handler = handlers.Logout(db)
	router = setupGetRouter(handler, route, "POST")
	writer = httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodPost, route, nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(wb, &jsonData); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	fmt.Printf("jsonData: %+v", jsonData)

	assert.True(t, jsonData["message"] == "successsfully logged out")
}

func TestUsers(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/users"
	handler := handlers.Users(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	// sampleUser ID
	marshalledData, err := json.Marshal("16")
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var res utils.ParsedUserRequestData
	if err := json.Unmarshal(wb, &res); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	assert.True(t, res.Email == "sample@gmail.com")
	assert.True(t, res.Username == "sampleUser")
	assert.True(t, res.Password == "")
}

func TestAuthenticateUser(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/login"
	handler := handlers.LoginUser(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	loginReq := utils.ParsedUserRequestData{
		Email:    "sample@gmail.com",
		Username: "sampleUser",
		Password: "sampleUser",
	}

	marshalledData, err := json.Marshal(loginReq)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	route = "/user"
	handler = handlers.AuthenticateUser(db)
	router = setupGetRouter(handler, route, "GET")

	newReq := httptest.NewRequest(http.MethodPost, route, nil)
	router.ServeHTTP(writer, newReq)

	assert.Equal(t, 200, writer.Code)
}

func TestLoginUser(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/login"
	handler := handlers.LoginUser(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	loginReq := utils.ParsedUserRequestData{
		Email:    "sample@gmail.com",
		Username: "sampleUser",
		Password: "sampleUser",
	}

	marshalledData, err := json.Marshal(loginReq)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	assert.Equal(t, 200, writer.Code)

	var res utils.ParsedUserRequestData
	if err := json.Unmarshal(wb, &res); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	assert.True(t, res.Email == "sample@gmail.com")
	assert.True(t, res.Username == "sampleUser")
	assert.True(t, res.Password == "")
}

func TestRegisterUser(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/sign-up"
	handler := handlers.RegisterUser(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	// type ParsedRequestData map[string]string
	signUpReq := utils.ParsedUserRequestData{
		Email:    "test@test.com",
		Username: "tester123",
		Password: "testerPassword",
	}

	marshalledData, err := json.Marshal(signUpReq)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	assert.Equal(t, 201, writer.Code)

	var res utils.ParsedUserRequestData
	if err := json.Unmarshal(wb, &res); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	assert.True(t, res.Email == "test@test.com")
	assert.True(t, res.Username == "tester123")
	assert.True(t, res.Password == "")

	var u models.Users
	db.Find(&u, "email = ?", res.Email)
	db.Unscoped().Delete(&u)
}

// tests fetching a true like and a nil like
func TestCheckArtworkLikes(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/likes"
	handler := handlers.CheckArtworkLikes(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	likeData := models.LikeReqData{
		ItemID:     "1015",
		UserID:     2,
		LikeStatus: true,
	}

	marshalledData, err := json.Marshal(likeData)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	assert.Equal(t, 202, writer.Code)

	var jsonData map[string]interface{}
	if err := json.Unmarshal(wb, &jsonData); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	assert.Equal(t, jsonData["liked"], true)
	assert.Equal(t, jsonData["exist"], true)

	// now we will test a false instance
	likeData = models.LikeReqData{
		ItemID:     "300",
		UserID:     1,
		LikeStatus: false,
	}

	marshalledDataFalse, err := json.Marshal(likeData)
	if err != nil {
		t.Error(err)
	}

	newWriter := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledDataFalse))
	router.ServeHTTP(newWriter, req)

	nwb, err := ioutil.ReadAll(newWriter.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var jsonDataFalse map[string]interface{}
	if err := json.Unmarshal(nwb, &jsonDataFalse); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonDataFalse: %s", err)
	}

	assert.Equal(t, jsonDataFalse["data"], nil)
	assert.Equal(t, jsonDataFalse["liked"], false)
	assert.Equal(t, jsonDataFalse["exist"], false)
}

// test will create a new instance of a like, then remove that like, then delete the instance
func TestArtworkLike(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/like"
	handler := handlers.ArtworkLike(db)
	router := setupGetRouter(handler, route, "POST")
	writer := httptest.NewRecorder()

	// inital like
	likeReq := models.LikeReqData{
		UserID:     2,
		ItemID:     "1000",
		LikeStatus: true,
	}

	marshalledData, err := json.Marshal(likeReq)
	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledData))
	router.ServeHTTP(writer, req)

	wb, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var like models.ArtworkLikes
	if err := json.Unmarshal(wb, &like); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to LikeReqData: %s", err)
	}

	assert.Equal(t, 201, writer.Code)
	assert.True(t, like.Like)

	// now to test unlike, IDforDeletion will be used for deleting at the end
	IDforDeletion := like.ID
	unlikeReq := models.LikeReqData{
		UserID:     2,
		ItemID:     "1000",
		LikeStatus: false,
	}

	marshalledDataUnlike, err := json.Marshal(unlikeReq)
	if err != nil {
		t.Error(err)
	}

	newWriter := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, route, bytes.NewReader(marshalledDataUnlike))
	router.ServeHTTP(newWriter, req)

	nwb, err := ioutil.ReadAll(newWriter.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(nwb, &jsonData); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to jsonData: %s", err)
	}

	assert.Equal(t, 200, newWriter.Code)
	assert.True(t, jsonData["message"] == "successfully updated")

	// delete instance
	var a models.ArtworkLikes
	db.Find(&a, "id = ?", IDforDeletion)
	db.Unscoped().Delete(&a)
}

func TestUsernames(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/usernames"
	handler := handlers.GetUsernames(db)
	router := setupGetRouter(handler, route, "GET")
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

	assert.True(t, true, len(usernames) >= 1)
}

func TestSearch(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/search/:term"
	handler := handlers.Search(db)
	router := setupGetRouter(handler, route, "GET")
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
}

func TestGetEra(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/era/:id"
	handler := handlers.GetEra(db)
	router := setupGetRouter(handler, route, "GET")
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
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/source/:id"
	handler := handlers.GetSource(db)
	router := setupGetRouter(handler, route, "GET")
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
}

func TestGetArtist(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/artist/:id"
	handler := handlers.GetArtist(db)
	router := setupGetRouter(handler, route, "GET")
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
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/artwork/:id"
	handler := handlers.GetArtwork(db)
	router := setupGetRouter(handler, route, "GET")
	writer := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/artwork/5", nil)
	router.ServeHTTP(writer, req)

	assert.Equal(t, 200, writer.Code)

	data, err := ioutil.ReadAll(writer.Body)
	if err != nil {
		t.Errorf("[Error] Unable to read writer.Body: %s", err)
	}

	var artwork models.Searches
	if err := json.Unmarshal(data, &artwork); err != nil {
		t.Errorf("[ERROR] Unable to unmarshal data to artwork: %s", err)
	}

	assert.Equal(t, "Portrait of the Cock Blomhoff family", artwork.Title)
	assert.Equal(t, "1817", artwork.DOR)
	assert.Equal(t, "https://lh6.ggpht.com/VK4zTxsKnaZnfOgNyJsMYtt-wf1aGV8rdpIlYCQs4azxhuo_Go3VXkAKR9INbOS5v5v2bREOnlQolXrmK6dsznV3VCw=s0", artwork.IMG)
}

func TestGetArtworks(t *testing.T) {
	db, _, err := utils.SetupConfiguration(true)
	if err != nil {
		t.Errorf("unable to setup db and env variables: %v", err)
	}

	route := "/artworks/"
	handler := handlers.GetArtworks(db)
	router := setupGetRouter(handler, route, "GET")
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
