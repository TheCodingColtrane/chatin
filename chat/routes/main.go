package routes

import (
	"messenger/controllers"

	"github.com/gorilla/mux"
)

func Start(router *mux.Router) {
	setChatRoutes(router)
}

func setChatRoutes(router *mux.Router) {
	router.HandleFunc("/chat", controllers.StartChat)
	router.HandleFunc("/chat/{id}", controllers.StartExistingChat)
}
