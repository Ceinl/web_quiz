package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quiz/handlers"
	"quiz/storage"
	"time"
)

func main() {
	db, err := storage.CreateDatabase("game.db")
	if err != nil {
    	log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/api/host", func(w http.ResponseWriter, r *http.Request) {
		handlers.HostHandler(db, w, r)
	})

	http.HandleFunc("/api/start_game", func(w http.ResponseWriter, r *http.Request) {
		handlers.Game(db, w, r)
	})

	http.HandleFunc("/api/connect", func(w http.ResponseWriter, r *http.Request) {
		handlers.Connect(db, w, r)
	})

	srv := &http.Server{Addr: ":8081"}
	
	go func() {
		log.Println("Backend running at httr://localhost:8081")
		if 	err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Printf("Server error: %v", err)
		}
			
	}()

	quit := make(chan os.Signal,1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Printf("Shutting down server...")
	ctx,canel := context.WithTimeout(context.Background(), 5*time.Second)
	defer canel()

	if err := srv.Shutdown(ctx); err != nil{
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Printf("Sever stopped")

} 
