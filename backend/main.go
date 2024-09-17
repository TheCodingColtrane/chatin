package main

import (
	"chatin/middleware"
	"chatin/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	var mux = mux.NewRouter()

	routes.Start(mux)
	http.ListenAndServe(":3000", middleware.EnableCORS(mux))
}
