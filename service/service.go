package service

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/efimovalex/wallet/adapters/database"
	"github.com/efimovalex/wallet/rest"
	"github.com/rs/zerolog/log"
	"github.com/synthesio/zconfig"
)

// Service - struct that aggregates and inits all the
// required packages for the service to run
type Service struct {
	database.Client `key:"database"`
	rest.REST       `key:"rest"`

	cancelContext context.Context
	errChan       chan error
}

// Run is a function to run the services
func (s *Service) Run(args []string) int {

	s.loadConfig()
	s.startService()

	return s.sigWait()
}

// load configuration data
func (s *Service) loadConfig() error {
	var err error // error holder
	zconfig.AddProviders(zconfig.Env)
	err = zconfig.Configure(s)
	if err != nil {
		log.Error().Err(err).Msg("error loading config")

		return err
	}

	// all okay
	return nil
}

// start the service listener
func (s *Service) startService() {
	s.errChan = make(chan error) // error response channel

	s.REST.DB = &s.Client
	// spawn REST listener
	go s.REST.Start(s.errChan)

}

// trap signals and wait
func (s *Service) sigWait() int {
	// trap signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// logging
	log.Info().Msg("Entering run loop")
	// run until signal or error
RunLoop:
	for {
		select {
		case sig := <-sigChan:
			// log signal
			log.Info().Msgf("Received signal: %d (%s)", sig, sig)

			if sig == syscall.SIGINT {

				err := s.REST.Stop()
				if err != nil {
					return 1
				}

				return 0
			}
			// check hup
			if sig == syscall.SIGHUP {
				// reload configuration
				_ = s.loadConfig()
				// log done
				log.Info().Msg("Configuration loaded: service changes may require restart")
				break
			}
			// break loop
			break RunLoop
		case err := <-s.errChan:
			// log error
			log.Error().Msgf("Received error from listener: %s", err.Error())
			return 1
		}
	}

	// exit clean
	return 0
}

// Synopsis shows the command summary
func (s *Service) Synopsis() string {
	return "Run the Wallet Service API"
}

// Help shows the detailed command options
func (s *Service) Help() string {

	zconfig.Configure(&Service{})
	zconfig.DefaultUsage(nil)

	return ""
}
