package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"real-time-forum/backend/api"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = utils.GeneratePort()
	}

	dbPath := "database/database.db"
	schemaPath := "database/schema.sql"
	db := database.NewDatabase(dbPath, schemaPath)
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("\033[31mError:\033[0m" + " Database shutdown failed - " + err.Error())
		} else {
			log.Println("\033[32mSuccess:\033[0m" + " Database connection closed gracefully")
		}
	}()

	wsHub := api.NewHub()
	go wsHub.StartHub()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: api.NewRouter(db, wsHub),
	}

	go func() {
		log.Printf("Server running on \033[36mhttp://localhost:%s\033[0m", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("\033[31mError:\033[0m" + " Server starting failed - " + err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wsHub.Shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("\033[31mError:\033[0m" + " Server shutdown failed - " + err.Error())
	}
	

	log.Println("\033[32mSuccess:\033[0m" + " Server shutdown gracefully")
}
