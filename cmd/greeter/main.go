package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"syscall"
	"time"

	"github.com/davidklassen/testing-workshop/pkg/dbutils"
	"github.com/davidklassen/testing-workshop/pkg/greetings"
	"github.com/davidklassen/testing-workshop/pkg/wait"
)

var (
	addr            = flag.String("addr", ":8080", "Application http network address")
	pgConnString    = flag.String("pg-conn-string", "postgres://greeter:greeter@pg:5432/greeter?sslmode=disable", "PostgreSQL server connection string")
	retryAttempts   = flag.Int("retry-attempts", 3, "Storage retry attempts")
	retryDelay      = flag.Duration("retry-delay", time.Second, "Storage retry delay")
	retryFactor     = flag.Float64("retry-factor", 2, "Storage retry backoff factor")
	shutdownTimeout = flag.Duration("shutdown-timeout", time.Second*30, "Graceful shutdown timeout")
)

func main() {
	log.Println("starting")

	db := dbutils.MustOpenPostgres(*pgConnString)
	greetingsRepo := greetings.NewRepo(db)

	greeter := NewGreeter(greetingsRepo, *retryAttempts, *retryDelay, *retryFactor)
	router := NewRouter(greeter)

	server := &http.Server{
		Addr:    *addr,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("http server failure: %s", err)
		}
	}()

	sig := wait.Wait(syscall.SIGTERM, syscall.SIGINT)
	log.Printf("received quit signal: %s, stopping", sig)

	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("http shutdown error: %s", err)
	}

	log.Println("bye-bye")
}
