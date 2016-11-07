package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type VersionNumber struct {
	Version float64 `json:"version"`
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
