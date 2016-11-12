package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/KevinMcIntyre/tu-nav-server/models"
	"github.com/julienschmidt/httprouter"
)

type VersionNumber struct {
	Version float64 `json:"version"`
}

type UpdateObject struct {
	Images    *[]Image           `json:"images"`
	Labels    *[]Label           `json:"labels"`
	Buildings *[]models.Building `json:"buildings"`
}

type Image struct {
	ID        int     `json:"id"`
	Name      *string `json:"name,omitempty"`
	Image     *string `json:"image,omitempty"`
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"long,omitempty"`
	Zoom      float64 `json:"zoom,omitempty"`
	Width     int     `json:"width,omitempty"`
	Height    int     `json:"height,omitempty"`
}

type Label struct {
	ID        int     `json:"id"`
	Name      *string `json:"name,omitempty"`
	Zoom      float64 `json:"zoom,omitempty"`
	Color     int     `json:"color,omitempty"`
	Size      int     `json:"size,omitempty"`
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"long,omitempty"`
}

func VerisonControlHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	returnResult := VersionNumber{getVersionNumber(DB)}
	fmt.Printf("returnResult.verison = %f\n", returnResult.Version)

	jsonResponse, err := json.Marshal(returnResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//TODO: Refactor error handling
	imageObjects, err := getImageData(DB)
	buildingObjects, err := models.GetBuildings(DB)
	labelObjects, err := getLabelData(DB)

	jsonObject := UpdateObject{imageObjects, labelObjects, buildingObjects}
	jsonResponse, err := json.Marshal(jsonObject)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}

func getVersionNumber(db *sql.DB) float64 {
	versionNumber := -1.0
	err := db.QueryRow(`SELECT VERSION FROM VERSION_NUMBER`).Scan(&versionNumber)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("VerisonControl - getVersionNumber - No row returned")
		break
	case err != nil:
		log.Fatal(err)
		break
	}
	return versionNumber
}

func getImageData(db *sql.DB) (*[]Image, error) {
	rows, err := db.Query(`SELECT
		id AS ID,
		name AS Name,
		image AS Image,
		lat AS Latitude,
		long AS Longitude,
    zoom AS Zoom,
    width AS Width,
    height AS Height
		FROM IMAGES`)
	if err != nil {
		return nil, err
	}
	var images []Image
	defer rows.Close()
	for rows.Next() {
		image := new(Image)
		if err := rows.Scan(
			&image.ID,
			&image.Name,
			&image.Image,
			&image.Latitude,
			&image.Longitude,
			&image.Zoom,
			&image.Width,
			&image.Height,
		); err != nil {
			return nil, err
		}
		images = append(images, *image)
	}
	return &images, nil
}

func getLabelData(db *sql.DB) (*[]Label, error) {
	rows, err := db.Query(`SELECT
		id AS ID,
		name AS Name,
    zoom AS Zoom,
    color AS Color,
    size AS Size,
    lat AS Latitude,
    long AS Longitude
		FROM LABELS`)
	if err != nil {
		return nil, err
	}
	var labels []Label
	defer rows.Close()
	for rows.Next() {
		label := new(Label)
		if err := rows.Scan(
			&label.ID,
			&label.Name,
			&label.Zoom,
			&label.Color,
			&label.Size,
			&label.Latitude,
			&label.Longitude,
		); err != nil {
			return nil, err
		}
		labels = append(labels, *label)
	}
	return &labels, nil
}
