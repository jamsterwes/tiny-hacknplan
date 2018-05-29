package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jamsterwes/tiny-hacknplan/server/api"

	"github.com/julienschmidt/httprouter"
)

// StartServer runs the server with the database at path `dbpath`.
func StartServer(apiKey string) {
	// http.ListenAndServe(":8080", http.FileServer(http.Dir("./public")))
	r := httprouter.New()
	r.GET("/", HomeHandler)
	r.GET("/assets/css/:stylesheet", CSSAssetHandler)
	r.GET("/assets/js/:script", JSAssetHandler)

	api.BootstrapAPI(&r, apiKey)

	http.ListenAndServe(":8080", r)
}

// HomeHandler handles templating the home page
func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	stylesheet, err := ioutil.ReadFile("./public/index.html")
	if err != nil {
		fmt.Fprintln(rw, err)
	}
	rw.Header().Add("Content-Type", "text/html")
	fmt.Fprintln(rw, string(stylesheet))
}

// CSSAssetHandler handles serving css files
func CSSAssetHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	filename := p.ByName("stylesheet")
	stylesheet, err := ioutil.ReadFile("./public/css/" + filename)
	if err != nil {
		fmt.Fprintln(rw, err)
	}
	rw.Header().Add("Content-Type", "text/css")
	fmt.Fprintln(rw, string(stylesheet))
}

// JSAssetHandler handles serving js files
func JSAssetHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	filename := p.ByName("script")
	script, err := ioutil.ReadFile("./public/js/" + filename)
	if err != nil {
		fmt.Fprintln(rw, err)
	}
	rw.Header().Add("Content-Type", "text/javascript")
	fmt.Fprintln(rw, string(script))
}
