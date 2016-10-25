package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KevinMcIntyre/tu-nav-server/models"
	"github.com/julienschmidt/httprouter"
)

func BuildingHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	buildings, err := models.GetBuildings(DB)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(buildings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}
