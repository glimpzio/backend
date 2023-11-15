package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/graph"
	"github.com/glimpzio/backend/profile"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	logger := log.New(os.Stdout, "[Gateway] ", log.Ldate|log.Ltime)

	if err := godotenv.Load(); err != nil {
		logger.Println(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		logger.Fatalln(err)
	}

	// Initialize services
	profileService := profile.NewProfileService(db)

	// Initialize handlers
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Logger: logger, ProfileService: profileService}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.ApplyMiddleware(srv))

	logger.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logger.Fatal(http.ListenAndServe(":"+port, nil))
}
