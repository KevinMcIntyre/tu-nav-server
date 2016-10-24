package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KevinMcIntyre/tu-nav-server/models"
	"github.com/julienschmidt/httprouter"
)

func ScheduleHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var scheduleRequest models.ScheduleRequest
	err := decoder.Decode(&scheduleRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scheduleValidtionIssues := scheduleRequest.Validate()
	if scheduleValidtionIssues != nil {
		jsonResponse, err := json.Marshal(scheduleValidtionIssues)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
		w.Write(jsonResponse)
		return
	}

	scheduleResponse, err := scheduleRequest.CallTUPortal()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(scheduleResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(jsonResponse)))
	w.Write(jsonResponse)
}
