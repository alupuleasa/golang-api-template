package main

import (
	"fmt"
	"os"

	"github.com/efimovalex/wallet/service"
	"github.com/mitchellh/cli"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// BuildName data
	BuildName = "wallet-service"
	// BuildDate data
	BuildDate string
	// BuildBranch data
	BuildBranch string
	// BuildNumber data
	BuildNumber string
)

// available commands
var cliCommands map[string]cli.CommandFactory

// init command factory
func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cliCommands = map[string]cli.CommandFactory{
		"run": func() (cli.Command, error) {
			return &service.Service{}, nil
		},
	}
}

func main() {
	var c *cli.CLI // cli object
	var status int // exit status
	var err error  // error holder

	// init and populate cli object
	c = cli.NewCLI(BuildName,
		fmt.Sprintf("%s-%s-%s",
			BuildBranch, BuildDate, BuildNumber),
	)
	c.Args = os.Args[1:]     // arguments minus command
	c.Commands = cliCommands // see commands above

	// run command and check return
	if status, err = c.Run(); err != nil {
		log.Error().Msgf("Error executing CLI: %s", err)
	}

	// exit
	os.Exit(status)
}
