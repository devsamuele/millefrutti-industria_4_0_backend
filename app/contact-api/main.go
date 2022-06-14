package main

// @title Elit Contact API documentation
// @version 1.0.0
// @host localhost:9000
// @BasePath /contact/v1

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/rs/cors"
)

func init() {
	time.Local = time.UTC
}

var build = "develop"

func main() {

	log := log.New(os.Stdout, "CONTACT-API: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	// Configuration
	var cfg struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:9000"`
			DebugHost       string        `conf:"default:0.0.0.0:8000"`
			ReadTimeout     time.Duration `conf:"default:10s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			ShutdownTimeout time.Duration `conf:"default:10s"`
		}
		Auth struct {
			Algorithm string `conf:"default:RS256"`
		}
		DB struct {
			URI     string        `conf:"default:mongodb+srv://samuele:s3b4n7c12@cluster0.7kbgj.mongodb.net/retryWrites=true&w=majority"`
			Name    string        `conf:"default:contact"`
			Timeout time.Duration `conf:"default:10s"`
		}
	}

	cfg.Version.SVN = build
	cfg.Version.Desc = "copyright info here"

	if err := conf.Parse(os.Args[1:], "CONTACT", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("CONTACT", &cfg)
			if err != nil {
				return fmt.Errorf("generating config usage: %w", err)
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("CONTACT", &cfg)
			if err != nil {
				return fmt.Errorf("generating config version: %w", err)
			}
			fmt.Println(version)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}

	log.Printf("main: Started: Application initializing: version %q", build)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Printf("main: Config:\n%v\n", out)

	// Authentication
	log.Println("main: Started: Initializing authentication support")

	// authentication, err := auth.New(cfg.Auth.Algorithm)
	// if err != nil {
	// 	return fmt.Errorf("construction auth: %w", err)
	// }

	// Database
	log.Println("main: Initializing database support")

	// Debug Service
	// /debug/pprof endpoint
	// /debug/vars endpoint
	expvar.NewString("build").Set(build)
	log.Println("main: Initializing debugging support")
	go func() {
		log.Printf("main: Debug Listening %s", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, http.DefaultServeMux); err != nil {
			log.Printf("main: Debug Listener closed : %v", err)
		}
	}()

	// Start API Service
	log.Println("main: Initializing API support")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           cors.AllowAll().Handler(handler.API(build, db, shutdown, log)),
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			_ = api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
