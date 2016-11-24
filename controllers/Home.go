package controllers

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/julienschmidt/httprouter"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var homeTemplate = pongo2.Must(pongo2.FromFile("templates/home.html"))
	err := homeTemplate.ExecuteWriter(pongo2.Context{"query": r.FormValue("query")}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
