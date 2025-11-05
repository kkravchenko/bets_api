package main

import (
	"context"
	"embed"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tz.api/internal/controller"
	"tz.api/internal/middleware"
	"tz.api/internal/store"
)

//go:embed *.html
var staticFiles embed.FS

func main() {
	port := flag.String("port", "8981", "HTTP port")
	flag.Parse()

	r := mux.NewRouter()

	var betStore = store.NewBetStore()
	betCtrl := controller.NewBetController(betStore)

	r.Use(middleware.LoggingMiddleware, middleware.RecoverMiddleware)

	// Роуты
	r.HandleFunc("/health", betCtrl.HealthHandler).Methods("GET")
	r.HandleFunc("/bets", betCtrl.CreateBetHandler).Methods("POST")
	r.HandleFunc("/bets", betCtrl.GetBetsHandler).Methods("GET")
	r.HandleFunc("/bets/{id}", betCtrl.GetBetByIDHandler).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.FS(staticFiles)))

	srv := &http.Server{
		Addr:    ":" + *port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on :%s", *port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}
