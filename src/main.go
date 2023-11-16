package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/graph"
	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func main() {
	logger := misc.NewLogger("Gateway", os.Stdout)

	if err := godotenv.Load(); err != nil {
		logger.ErrorLog.Println(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.ErrorLog.Fatalln(err)
	}

	auth0Config := &auth.Auth0Config{
		Auth0Domain:       os.Getenv("AUTH0_DOMAIN"),
		Auth0ClientId:     os.Getenv("AUTH0_CLIENT_ID"),
		Auth0ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
	}

	mailList := misc.NewMailList(os.Getenv("SENDGRID_API_KEY"), os.Getenv("SENGRID_LIST_ID"))

	// Initialize services
	profileService := profile.NewProfileService(db, mailList)

	// Initialize handlers
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Logger: logger, ProfileService: profileService}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", auth.ApplyMiddleware(logger, srv, auth0Config))

	logger.InfoLog.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logger.ErrorLog.Fatal(http.ListenAndServe(":"+port, nil))
}
