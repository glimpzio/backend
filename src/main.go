package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/glimpzio/backend/graph"
	"github.com/joho/godotenv"
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

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
