package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dbubel/bchat/cmd/bchat-api/handlers"
	"github.com/dbubel/bchat/internal/platform/db"

	"github.com/kelseyhightower/envconfig"
)

var build = "develop"

func main() {
	log := log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// API Server config
	var cfg struct {
		Web struct {
			APIHost         string        `default:":3000" envconfig:"API_HOST"`
			ReadTimeout     time.Duration `default:"5s" envconfig:"READ_TIMEOUT"`
			WriteTimeout    time.Duration `default:"5s" envconfig:"WRITE_TIMEOUT"`
			ShutdownTimeout time.Duration `default:"5s" envconfig:"SHUTDOWN_TIMEOUT"`
		}
		DB struct {
			Host string `default:"bfscans.db" envconfig:"DB_HOST"`
		}
	}

	if err := envconfig.Process("TRACE-API", &cfg); err != nil {
		log.Fatalf("Parsing Config : %v", err)
	}

	// =========================================================================
	// Start MySQL

	log.Println("main : Started : Initialize SQLite")
	var err error
	var masterDB *db.SQLite
	for i := 0; i < 10; i++ {

		masterDB, err = db.New(cfg.DB.Host)
		if err != nil {
			log.Printf("main : Register DB : %s\n", err.Error())
		} else {
			log.Println("DB connect OK")
			break
		}
		time.Sleep(time.Second)
	}

	// =========================================================================
	// Start API
	api := http.Server{
		Addr:           cfg.Web.APIHost,
		Handler:        handlers.API(log, masterDB),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    cfg.Web.ReadTimeout,
		WriteTimeout:   cfg.Web.WriteTimeout,
	}

	log.Printf("Started : Application Initializing version %q", build)
	defer log.Println("App Shutdown")

	cfgJSON, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatalf("Marshalling Config to JSON : %v", err)
	}
	log.Printf("%v\n", string(cfgJSON))

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("API Listening %s", cfg.Web.APIHost)
		serverErrors <- api.ListenAndServe()
	}()

	// Shutdown

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Stop API Service
	// Blocking main and waiting for graceful shutdown and connection draining.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)

	case <-osSignals:
		log.Println("Start shutdown...")

		// Create context for Shutdown call.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		if err := api.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			if err := api.Close(); err != nil {
				log.Fatalf("Could not stop http server: %v", err)
			}
		}
	}
}
