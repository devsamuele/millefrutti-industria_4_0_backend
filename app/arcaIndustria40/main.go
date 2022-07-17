package main

// @title Elit Contact API documentation
// @version 1.0.0
// @host localhost:9000
// @BasePath /contact/v1

import (
	"context"
	"encoding/json"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/devsamuele/millefrutti-industria_4_0_backend/app/arcaIndustria40/handler"
	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/data/pasteurizer"
	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/data/spindryer"
	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/sys/database"
	"github.com/devsamuele/service-kit/ws"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/rs/cors"
)

// func init() {
// 	time.Local = time.UTC
// }

var build = "develop"

func main() {

	log := log.New(os.Stdout, "ARCA-INDUSTRIA-4-0-API: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
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
			// URI     string        `conf:"default:localhost"`
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

	// Database
	log.Println("main: Initializing database support")
	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("main: opening db: %w", err)
	}

	go func() {
		var connected bool
		for {
			if err := db.Ping(); err != nil {
				if connected {
					log.Println("sql: db", err)
				}
				connected = false
			} else {
				if !connected {
					log.Printf("sql: db connected")
				}
				connected = true
			}
			time.Sleep(time.Second * 5)
		}
	}()

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

	// websocket init
	log.Println("main: Initializing websocket support")
	io := ws.New(nil)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// OPCUA Spindryer Service
	log.Println("main: Initializing opcua support")
	ctx := context.Background()

	// Pasteurizer "opc.tcp://192.168.1.201:4840"
	pasteurizerClient := opcua.NewClient("opc.tcp://192.168.1.201:4840", opcua.SecurityMode(ua.MessageSecurityModeNone), opcua.DialTimeout(time.Second*10))
	defer pasteurizerClient.CloseWithContext(ctx)

	// Spindryer "opc.tcp://192.168.1.22:4840"
	spindryerClient := opcua.NewClient("opc.tcp://192.168.1.22:4840", opcua.SecurityMode(ua.MessageSecurityModeNone), opcua.DialTimeout(time.Second*10))
	defer spindryerClient.CloseWithContext(ctx)

	// TODO Move away
	go func() {
		if spindryerClient.State() == opcua.Connected {
			spindryer.OpcuaConnected = true
			opcuaSpindryerService := spindryer.NewOpcuaService(ctx, log, spindryerClient, shutdown, spindryer.NewStore(db, log), &io)
			opcuaSpindryerService.Run()
			conn := spindryer.OpcuaConnection{
				Connected: pasteurizer.OpcuaConnected,
			}
			b, err := json.Marshal(conn)
			if err != nil {
				log.Println(err)
			}
			if err := io.Broadcast("spindryer-opcua-connection", b); err != nil {
				log.Println(err)
			}
		}

		for {
			if spindryerClient.State() != opcua.Connected {
				if err := spindryerClient.Connect(ctx); err != nil {
					if spindryer.OpcuaConnected {
						log.Printf("opcua: spindryer: %v", err)
					}
					spindryer.OpcuaConnected = false
				} else {
					if !spindryer.OpcuaConnected {
						log.Printf("opcua: spindryer connected")
					}
					spindryer.OpcuaConnected = true
				}

				conn := spindryer.OpcuaConnection{
					Connected: pasteurizer.OpcuaConnected,
				}
				b, err := json.Marshal(conn)
				if err != nil {
					log.Println(err)
				}
				if err := io.Broadcast("spindryer-opcua-connection", b); err != nil {
					log.Println(err)
				}
			}
			time.Sleep(time.Second * 15)
		}
	}()

	go func() {
		if pasteurizerClient.State() == opcua.Connected {
			pasteurizer.OpcuaConnected = true
			opcuaPasteurizerService := pasteurizer.NewOpcuaService(ctx, log, pasteurizerClient, shutdown, pasteurizer.NewStore(db, log), &io)
			opcuaPasteurizerService.Run()
			conn := pasteurizer.OpcuaConnection{
				Connected: pasteurizer.OpcuaConnected,
			}
			b, err := json.Marshal(conn)
			if err != nil {
				log.Println(err)
			}
			if err := io.Broadcast("pasteurizer-opcua-connection", b); err != nil {
				log.Println(err)
			}
		}

		for {
			if pasteurizerClient.State() != opcua.Connected {
				if err := pasteurizerClient.Connect(ctx); err != nil {
					if pasteurizer.OpcuaConnected {
						log.Printf("opcua: pasteurizer: %v", err)
					}
					pasteurizer.OpcuaConnected = false
				} else {
					if !pasteurizer.OpcuaConnected {
						log.Printf("opcua: pasteurizer connected")
					}
					pasteurizer.OpcuaConnected = true
				}
				conn := pasteurizer.OpcuaConnection{
					Connected: pasteurizer.OpcuaConnected,
				}
				b, err := json.Marshal(conn)
				if err != nil {
					log.Println(err)
				}
				if err := io.Broadcast("pasteurizer-opcua-connection", b); err != nil {
					log.Println(err)
				}
			}

			time.Sleep(time.Second * 15)
		}
	}()

	// Start API Service
	log.Println("main: Initializing API support")

	serverErrors := make(chan error, 1)

	api := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           cors.AllowAll().Handler(handler.API(build, db, spindryerClient, pasteurizerClient, &io, shutdown, log)),
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

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
