package handler

import (
	"github.com/gorilla/mux"
)

// InitRouter инициализирует роутер для обработчиков HTTP-запросов.
func InitRouter() *mux.Router {
	router := mux.NewRouter()

	// Регистрация обработчиков.
	router.HandleFunc("/api/update", updateHandler).Methods("POST")
	router.HandleFunc("/api/state", stateHandler).Methods("GET")
	router.HandleFunc("/api/get_names_strong", getNamesStrong).Methods("GET")
	router.HandleFunc("/api/get_names_weak", getNamesWeak).Methods("GET")
	router.HandleFunc("/api/get_names", getNames).Methods("GET")

	return router
}
