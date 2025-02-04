package main

import (
	"app/backend/api"
	"app/backend/db"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	db, err := db.NewDatabase("./db/forum.db", "./db/schema.sql", db.Config{
		MaxOpenConns:    10, // max number of open connections to the database
		MaxIdleConns:    5,  // max number of idle connections in the connection pool
		ConnMaxLifetime: 10 * time.Minute, // maximum amount of time a connection can be reused
	})
	if err != nil {
		log.Fatal(err)

	}
	port := os.Getenv("PORT")
	if port == "" {
		port = findFreePort()
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: api.NewRouter(db),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Printf("\nServer running on http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERROR: Server failed - %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("\nShutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// closing database connection
	if err := db.Close(); err != nil {
		log.Printf("ERROR: Failed to close database connection - %v", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("ERROR: Shutdown failed - %v", err)
	}

	wg.Wait()
	log.Println("Clean shutdown complete")
}

func findFreePort() string {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("ERROR: Failed to find free port - %v", err)
	}
	defer listener.Close()
	return fmt.Sprintf("%d", listener.Addr().(*net.TCPAddr).Port)
}
