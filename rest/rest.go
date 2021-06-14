package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// REST - struct that aggregates and inits all the
// required packages for the REST API to run
type REST struct {
	router *httprouter.Router
	server *http.Server

	DB Database

	Addr string `key:"addr" description:"address the server should bind to" default:":80"`
}

type Database interface{}

func (r *REST) Init() (err error) {
	return nil
}

// Start - Starts the http listener
func (r *REST) Start(e chan error) {
	r.router = httprouter.New()

	r.loadRoutes()

	// log startup
	log.Info().Msg(fmt.Sprintf("HTTP Service listening on %s", r.Addr))

	r.server = &http.Server{
		Addr:         r.Addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      r.router,
	}
	// send feedback on listener error
	e <- r.server.ListenAndServe()
}

func (r *REST) Stop() error {
	var err error // error holder
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = r.server.Shutdown(ctxShutDown); err != nil {
		log.Error().Err(err).Msg("server Shutdown Failed")

		return err
	}

	log.Info().Msg("Server stopped properly")
	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}
