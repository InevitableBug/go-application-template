package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func initialiseRoutes() http.Handler {

	router := mux.NewRouter()

	router.Use(handlers.RecoveryHandler())

	corsOrigins := handlers.AllowedOrigins([]string{"*"})
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	router.Use(handlers.CORS(corsMethods, corsOrigins))

	//------------------------------------------------------------------------------------------------------------------
	// status endpoint
	//------------------------------------------------------------------------------------------------------------------

	statusRouter := router.PathPrefix("/_status").Subrouter()

	statusRouter.HandleFunc("/ping", pingHandler)

	//------------------------------------------------------------------------------------------------------------------
	// api endpoints
	//------------------------------------------------------------------------------------------------------------------
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(responseContentTypeMiddleware)

	//------------------------------------------------------------------------------------------------------------------
	// static file server
	//------------------------------------------------------------------------------------------------------------------

	fs := http.Dir(config.HTTP.DocumentRoot)

	router.PathPrefix("/").Handler(http.FileServer(fs))

	return router
}

func responseContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}
