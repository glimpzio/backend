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

	auth0Config := &auth.Auth0Config{
		Auth0Domain:       os.Getenv("AUTH0_DOMAIN"),
		Auth0ClientId:     os.Getenv("AUTH0_CLIENT_ID"),
		Auth0ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
	}

	// Initialize services
	profileService := profile.NewProfileService(db)

	// Initialize handlers
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Logger: logger, ProfileService: profileService}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.ApplyMiddleware(logger, srv, auth0Config))

	logger.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logger.Fatal(http.ListenAndServe(":"+port, nil))
}
