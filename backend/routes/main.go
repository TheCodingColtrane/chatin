package routes

import (
	"chatin/controllers"
	"chatin/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

const GET = "GET"
const POST = "POST"

func Start(router *mux.Router) {
	userRoutes(router)
	chatRoutes(router)
	authRoutes(router)
	searchRoutes(router)
}

func userRoutes(router *mux.Router) {

	var controller = controllers.NewUserController()
	var userRouter = router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", controller.Get).Methods(GET)
	userRouter.HandleFunc("", controller.Create).Methods(POST)
	userRouter.HandleFunc("/data", authenticate(controller.Find)).Methods(GET)

}

func searchRoutes(router *mux.Router) {
	var controller = controllers.NewSearchController()
	router.HandleFunc("/search", authenticate(controller.GetSearchFoundItems)).Methods(POST)
	router.HandleFunc("/search/items-count", authenticate(controller.GetSearchFoundItemsNumber)).Methods(POST)

}

func chatRoutes(router *mux.Router) {
	var controller = controllers.NewChatsController()
	var chatRouter = router.PathPrefix("/chats").Subrouter()
	chatRouter.HandleFunc("/", authenticate(controller.GetChatList)).Methods(GET)
	chatRouter.HandleFunc("/{chatID}", authenticate(controller.Get)).Methods(GET)
	chatRouter.HandleFunc("/{chatID}/messages", authenticate(controller.GetMoreMessages)).Methods(GET)

}

func authRoutes(router *mux.Router) {
	var controller = controllers.NewAuthController()
	var authRouter = router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("", controller.Login)
}

func authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return middleware.Authenticate(handler).ServeHTTP
}
