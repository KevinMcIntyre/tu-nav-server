package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KevinMcIntyre/tu-nav-server/utils"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

func HelloHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response, err := json.Marshal(struct {
		Payload string `json:"payload"`
	}{"Hello world from TU Nav!"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	w.Write(response)
}

func main() {
	utils.WritePid()

	router := httprouter.New()
	router.GET("/hello", HelloHandler)

	n := negroni.New(
		negroni.NewRecovery(),
	)

	n.UseHandler(router)
	n.Run(":3030")
}
