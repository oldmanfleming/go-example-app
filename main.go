package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/eflem00/go-example-app/controllers"
	"github.com/eflem00/go-example-app/controllers/http"
	"github.com/eflem00/go-example-app/controllers/queue"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
}

func configLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env := os.Getenv("ENV"); env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func awaitSigterm() {
	log.Info().Msg("awaiting sigterm")

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	sig := <-cancelChan

	log.Info().Msgf("caught sigterm %v", sig)
}

func main() {

	loadEnv()

	configLogger()

	log.Info().Msg("starting app")

	// start a slice of blocking functions in concurrent go routines
	// functions implement IController
	contrs := []controllers.IController{
		http.HttpController{},
		queue.QueueController{},
	}

	for _, contr := range contrs {
		go func(contr controllers.IController) {
			// start is intended to be a blocking call
			// if Exit() is called, we have caught a panic
			// if start returns, one of our controllers is no longer active and thus we should force a panic
			defer contr.Exit()
			err := contr.Start()
			panic(err)
		}(contr)
	}

	// blocking call in main routine to await sigterm
	awaitSigterm()

	// TODO: Shutdown gracefully below
}
