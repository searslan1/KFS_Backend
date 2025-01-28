package user

import "github.com/gorilla/mux"

func RegisterUserRoutes(router *mux.Router, controller *UserController) {
	router.HandleFunc("/users", controller.RegisterHandler).Methods("POST")
	router.HandleFunc("/users/login", controller.LoginHandler).Methods("POST")
}
