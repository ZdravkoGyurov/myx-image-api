package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/api/router"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/config"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/controller"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/log"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/storage/postgresql"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/storage/s3"

	"github.com/rs/cors"
)

type GlobalContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}

func NewGlobalContext() GlobalContext {
	ctx, cancel := context.WithCancel(context.Background())
	return GlobalContext{
		Context: ctx,
		Cancel:  cancel,
	}
}

type Application struct {
	appContext GlobalContext
	config     config.Config
	server     *http.Server
	psql       *postgresql.Storage
	s3         *s3.Storage
}

func New(cfg config.Config) *Application {
	globalContext := NewGlobalContext()

	psql := postgresql.New(cfg.PostgreSQL)
	s3 := s3.New(cfg.S3)
	ctrl := controller.New(cfg, psql, s3)
	r := router.New(ctrl)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	enableCors(r, server)

	return &Application{
		appContext: globalContext,
		config:     cfg,
		server:     server,
		psql:       psql,
		s3:         s3,
	}
}

func enableCors(router router.Router, server *http.Server) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5500"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodPut,
			http.MethodOptions,
			http.MethodDelete,
		},
		Debug: true,
	})
	server.Handler = c.Handler(router)
}

func (a *Application) Start() error {
	logger := log.DefaultLogger()
	logger.Info().Msg("starting application...")
	a.setupSignalNotifier()

	if err := a.psql.Connect(a.appContext.Context); err != nil {
		return err
	}
	logger.Info().Msg("postgresql connection opened")

	if err := a.s3.Connect(); err != nil {
		return err
	}
	logger.Info().Msg("s3 connection opened")

	logger.Info().Msg("server started")
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Err(err).Msg("failed listen and serve")
		}
	}()

	logger.Info().Msg("application started")

	<-a.appContext.Context.Done()

	a.stopServer()
	a.closePSQLConnection()
	logger.Info().Msg("application stopped")
	return nil
}

func (a *Application) Stop() {
	a.appContext.Cancel()
}

func (a *Application) setupSignalNotifier() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.DefaultLogger().Info().Msg("stopping application...")
		a.appContext.Cancel()
	}()
}

func (a *Application) closePSQLConnection() {
	a.psql.Close()
	log.DefaultLogger().Info().Msg("storage connection closed")
}

func (a *Application) stopServer() {
	logger := log.DefaultLogger()
	ctx, cancel := context.WithTimeout(context.Background(), a.config.Server.ShutdownTimeout)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Err(err).Msg("failed to stop server")
	}
	logger.Info().Msg("server stopped")
}
