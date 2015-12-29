package handlers

import (
	"net/http"

	"github.com/gorilla/pat"
)

func GetRouter() *pat.Router {
	r := pat.New()

	r.Get("/v1/health", HealthHandler)
	r.Post("/v1/register", RegisterCreateHandler)
	r.Delete("/v1/register", RegisterDeleteHandler)

	//Add routing for static routes
	s := http.StripPrefix("/swagger/", http.FileServer(http.Dir("/swagger")))
	r.PathPrefix("/swagger/").Handler(s)

	return r
}
