package main

import (
	"context"
	"cryptocurrency/internal/collectors"
	"cryptocurrency/internal/observers"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Cryptocurrency Bot API documentation
// @version 0.1.0

// @host localhost:8080
// @BasePath /api

//@securityDefinitions.apikey Bearer
//@in header
//@name Authorization

func main() {
	// Load env variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading env file")
	}
	log.Println("Env variables loaded")

	// Initialize Data sources
	ds, err := initDS()
	if err != nil {
		log.Fatalf("Unable to initialize data sources: %v\n", err)
	}

	// Init cron tasks
	cronObservers, err := initCronTasks(ds)
	if err != nil {
		log.Fatalf("Failure to create cron tasks: %v\n", err)
	}

	// Inject
	router, err := inject(ds, cronObservers)
	if err != nil {
		log.Fatalf("Failure to inject data sources: %v\n", err)
	}

	// Init general data about currency
	err = initGeneralData(ds)
	if err != nil {
		log.Fatalf("Failure to initialize general data: %v\n", err)
	}

	// Initialize Cron tasks
	c := cron.New()
	// collect prices
	_, err = c.AddFunc("@every 10s", func() {
		collectors.GetPrice(ds.BadgerClient, ds.MongoDBClient)
	})
	_, err = c.AddFunc("@every 20m", func() {
		observers.MakeObserverForPercentage(ds.BadgerClient, ds.MongoDBClient, time.Minute*20)
	})
	_, err = c.AddFunc("@every 1h", func() {
		observers.MakeObserverForPercentage(ds.BadgerClient, ds.MongoDBClient, time.Minute*60)
	})
	if err != nil {
		log.Fatalf("Failure to start cron tasks: %v\n", err)
	}
	c.Start()
	log.Println("Cron task started")

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Graceful server shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("Listening on port %v\n", server.Addr)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish the request it's currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown data source
	if err := ds.closeDS(); err != nil {
		log.Fatalf("A problem occured gracefully shutting down data sources: %v\n", err)
	}

	// Shutdown server
	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	// Stop the scheduler
	c.Stop()
}
