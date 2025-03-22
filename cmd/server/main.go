package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/h1sashin/go-app/config"
	"github.com/h1sashin/go-app/db"
	translator "github.com/h1sashin/go-app/i18n"
	"github.com/h1sashin/go-app/logging"
	"github.com/h1sashin/go-app/middleware"
	"github.com/h1sashin/go-app/service"
	"github.com/rs/zerolog/log"

	// admin "github.com/h1sashin/go-app/graph/admin/generated"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	public "github.com/h1sashin/go-app/graph/public/generated"
	publicResolver "github.com/h1sashin/go-app/graph/public/resolver"
)

func main() {
	logging.SetupLogger(config.LogLevel(config.Info))

	cfg, err := config.Load()

	if err != nil {
		log.Fatal().Err(err)
	}

	conn, err := db.NewDB(cfg)

	if err != nil {
		log.Fatal().Err(err)
	}

	defer conn.Close(context.Background())

	bundle := translator.NewTranslator()

	queries := db.New(conn)

	service := service.NewService(conn, queries, cfg)

	router := http.NewServeMux()

	public := handler.New(public.NewExecutableSchema(public.Config{Resolvers: &publicResolver.Resolver{Cfg: cfg, Service: service}}))
	public.AddTransport(transport.Options{})
	public.AddTransport(transport.POST{})

	if cfg.AppEnv == config.Development {
		public.Use(extension.Introspection{})
		router.Handle("/graphql/playground", playground.Handler("GraphQL playground", "/graphql"))
	}

	router.Handle("/graphql", public)

	stack := middleware.CreateStack(middleware.Localizer(bundle))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.AppPort),
		Handler: stack(router),
	}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("Shutting down server...")
		os.Exit(0)
	}()

	log.Printf("Server running on port %d", cfg.AppPort)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err)
	}

}
