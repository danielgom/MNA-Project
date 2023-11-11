// Package main where the execution of the program lives.
package main

import (
	"MNA-project/pkg/config/core"
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"MNA-project/pkg/config"
	"MNA-project/pkg/config/logger"

	_ "github.com/swaggo/echo-swagger/example/docs"
)

const defaultServerTimeout = time.Second * 5

// @title PetSys
// @version 1.0
// @description This is the API server for PetSys application.
// @termsOfService http://swagger.io/terms/

// @contact.name Daniel G_A
// @contact.url https://tec.mx/es
// @contact.email A01794498@tec.mx

// @host localhost:8080
func main() {
	logger.Initialise()
	r := core.Router()

	r.Server = &http.Server{
		ReadTimeout:       defaultServerTimeout,
		WriteTimeout:      defaultServerTimeout,
		IdleTimeout:       defaultServerTimeout,
		ReadHeaderTimeout: defaultServerTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Start(config.LoadConfig().Port); err != nil {
			log.Fatalln("error serving", err.Error())
		}
	}()

	<-ctx.Done()

	stop()

	log.Println("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), defaultServerTimeout)
	if err := r.Shutdown(ctx); err != nil {
		cancel()
		log.Fatalln("server forced to shutdown: ", err)
	}
}
