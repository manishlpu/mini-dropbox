package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/manishlpu/assignment/api"
	"github.com/manishlpu/assignment/utils"
	"github.com/joho/godotenv"
)

func NewServer() (*http.Server, error) {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	api, err := api.New()
	if err != nil {
		return nil, err
	}

	srvHost := utils.GetEnvValue("APP_HOST", "localhost")
	srvPort := utils.GetEnvValue("APP_PORT", "8081")
	srvAddress := fmt.Sprintf("%s:%v", srvHost, srvPort)
	log.Println("Configuring Server at address ", srvAddress)
	srv := http.Server{
		Addr:    srvAddress,
		Handler: api,
		// Read will Timeout after 2s if anything goes wrong.
		ReadTimeout: time.Duration(2 * time.Second),
	}

	return &srv, nil
}

func StartServer(srv *http.Server) {
	log.Println("Starting Server...")

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down the server gracefully...")
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Println("HTTP server Shutdown: ", err)
			return
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Println("HTTP server ListenAndServe: ", err)
		return
	}

	<-idleConnsClosed
}
