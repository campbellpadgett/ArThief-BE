package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type Era struct {
	ID            int       `json:"id"`
	Era_Name      string    `json:"era_name"`
	Last_Modified time.Time `json:"last_modified"`
}

func (Era) TableName() string {
	return "artwork_migrate_era"
}

type Artwork struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Nationality     string    `json:"nationality"`
	Artist_Bio      string    `json:"artist_bio"`
	Desc            string    `json:"desc"`
	Culture         string    `json:"culture"`
	Gender          string    `json:"gender"`
	Nation          string    `json:"nation"`
	Medium          string    `json:"medium"`
	Date_of_Release string    `json:"date_of_release"`
	Image           string    `json:"image"`
	Image_Small     string    `json:"image_small"`
	Last_Modified   time.Time `json:"last_modified"`
	Artist_ID       int       `json:"artist_id"`
	Source_ID       int       `json:"source_id"`
}

func (Artwork) TableName() string {
	return "artwork_migrate_artwork"
}

type Artist struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Era           string    `json:"era"`
	Gender        string    `json:"gender"`
	Last_Modified time.Time `json:"last_modified"`
}

func (Artist) TableName() string {
	return "artwork_migrate_artist"
}

type Source struct {
	ID            int       `json:"id"`
	Source_Name   string    `json:"source_name"`
	Abriviation   string    `json:"abriviation"`
	Last_Modified time.Time `json:"last_modified"`
}

func (Source) TableName() string {
	return "artwork_migrate_source"
}

type Searches struct {
	ID          string `json:"id"`
	Title       string `json:"Title"`
	Artist_Name string `json:"Artist_Name"`
	DOR         string `json:"DOR"`
	Description string `json:"Description"`
	Source      string `json:"Source"`
	Abb         string `json:"Abb"`
	IMG         string `json:"IMG"`
	IMG_S       string `json:"IMG_S"`
}

func (Searches) TableName() string {
	return "searches"
}

type Users struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
}

func (Users) TableName() string {
	return "users"
}

// Curations are created when a user supplies thier userID and a name. When adding
// artwork, we persist the artworkID to CurationArtwork table. We join the tables
// together, along with the searches table and write the artwork objects to Curations.Artworks
type Curations struct {
	gorm.Model
	User_ID  int       `json:"user_id"`
	Name     string    `json:"name"`
	Artworks []Artwork `json:"artworks" gorm:"many2many:curation_artwork"`
}

func (Curations) TableName() string {
	return "curations"
}

type CurationLikes struct {
	gorm.Model
	Curation_ID uint `json:"curation_id"`
	User_ID     int  `json:"user_id"`
	Like        bool `json:"like"`
}

func (CurationLikes) TableName() string {
	return "curation_likes"
}

type CurationArtwork struct {
	gorm.Model
	Curation_ID int `json:"curation_id"`
	Artwork_ID  int `json:"user_id"`
	Order       int `json:"order"`
}

func (CurationArtwork) TableName() string {
	return "curation_artwork"
}

type ArtworkLikes struct {
	gorm.Model
	Artwork_ID int  `json:"artwork_id"`
	User_ID    int  `json:"user_id"`
	Like       bool `json:"like"`
}

func (ArtworkLikes) TableName() string {
	return "artwork_likes"
}

type LikeReqData struct {
	ItemID     string
	LikeStatus bool
	UserID     int
}

// Returns string of LikeReqData
func (d *LikeReqData) ToString() string {
	return fmt.Sprintf("IID: %v, UID: %v, L: %v", d.ItemID, d.UserID, d.LikeStatus)
}

// Takes in request and processes the body for an instance of LikeReqData
func (d *LikeReqData) ProcessReq(req *http.Request) error {
	data, ioErr := ioutil.ReadAll(req.Body)
	if ioErr != nil {
		return ioErr
	}

	if mErr := json.Unmarshal(data, &d); mErr != nil {
		return mErr
	}

	return nil
}

type LikedList struct {
	LikedArtwork []Searches `json:"liked_artwork"`
	NextPage     int        `json:"page"`
	Count        int64      `json:"count"`
}

func (ll *LikedList) AddNextPage(amt int) (int, error) {
	if amt <= 0 {
		return 0, errors.New("amt param cannot be less than or equal to 0")
	}

	oldValue := ll.NextPage
	ll.NextPage += amt

	if ll.NextPage == oldValue {
		return 0, errors.Wrapf(errors.New("NextPage did not increase"), "old: %v, new: %v", oldValue, ll.NextPage)
	}

	return ll.NextPage, nil
}
