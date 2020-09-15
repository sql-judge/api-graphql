package main

import (
	"context"
	"flag"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/sql-judge/api-graphql/config"
	"github.com/sql-judge/api-graphql/database"
	"github.com/sql-judge/api-graphql/graph"
	"github.com/sql-judge/api-graphql/graph/generated"
	"log"
	"net/http"
)

const defaultConfigPath = "config.yml"

func main() {
	// set up command-line flags
	var configPath string
	flag.StringVar(&configPath, "cfg", defaultConfigPath, "configuration file path")
	flag.Parse()

	// load configuration
	cfg := config.Config{}
	if err := cfg.LoadFromFile(configPath); err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	// connect to the database
	db, err := database.ConnectWithConfig(cfg.DatabaseConfig)
	if err != nil {
		log.Fatalf("failed to connect to the database: %s", err)
	}
	defer db.Close(context.Background())

	// create router
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:3000"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	// create and start server
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(db)}))

	router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	router.Handle("/query", srv)

	err = http.ListenAndServe(cfg.ServerConfig.Address(), router)
	if err != nil {
		panic(err)
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cfg.ServerConfig.Port)
}
