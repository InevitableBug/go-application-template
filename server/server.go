package server

import (
	"fmt"
	"net/http"

	bugConfig "bug.geek.nz/go-application-template/config"
	log "github.com/sirupsen/logrus"
)

var config = bugConfig.Instance

func Start() *http.Server {
	address := fmt.Sprintf(":%d", config.HTTP.Port)
	log.Info(fmt.Sprintf("Initialising web service.  Listening on address '%s'", address))

	handler := initialiseRoutes()

	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Error(err)
		}
	}()

	return server
}
