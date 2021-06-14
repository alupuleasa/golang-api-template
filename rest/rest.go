package rest

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/efimovalex/wallet/adapters/model"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

// REST - struct that aggregates and inits all the
// required packages for the REST API to run
type REST struct {
	router *httprouter.Router
	server *http.Server

	DB Database

	Addr    string `key:"addr" description:"address the server should bind to" default:":80"`
	AuthKey string
}

// Database - interface with the database
type Database interface {
	CreateWallet(OwnerAccountID uint64) (w *model.Wallet, err error)
	GetWallets(limit, offset uint64) (wallets []model.Wallet, err error)
	RemoveFunds(walletID uint64, sum float64) (err error)
	AddFunds(walletID uint64, sum float64) (err error)
}

// Start - Starts the http listener
func (r *REST) Start(e chan error) {
	// simulate an auth system
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		e <- err

		return
	}

	r.AuthKey = fmt.Sprintf("%x", key)
	log.Info().Msgf("Auth API KEY: %s", r.AuthKey)

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

// Stop - gracefully stops the REST API
func (r *REST) Stop() error {
	var err error // error holder
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer func() {
		cancel()
	}()

	if err = r.server.Shutdown(ctxShutDown); err != nil {
		log.Error().Err(err).Msg("server shutdown failed")

		return err
	}

	log.Info().Msg("Server stopped properly")
	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}

// WriteError - formats the error response
func (r *REST) WriteError(w http.ResponseWriter, statusCode int, err string) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)

	errStruct := struct {
		ErrorStatus     string `json:"status,omitempty"`
		ErrorStatusCode int    `json:"code,omitempty"`
		Message         string `json:"message,omitempty"`
	}{
		http.StatusText(statusCode),
		statusCode,
		err,
	}
	json.NewEncoder(w).Encode(errStruct)
}
