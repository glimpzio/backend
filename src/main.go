package main

import (
	"database/sql"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/graph"
	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const defaultPort = "8080"

func graphqlHandler(logger *misc.Logger, auth0Config *auth.Auth0Config, resolver *graph.Resolver) gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	wrapped := auth.ApplyMiddleware(logger, h, auth0Config)

	return func(c *gin.Context) {
		wrapped.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

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
	r := gin.Default()
	r.POST("/query", graphqlHandler(logger, auth0Config, &graph.Resolver{Logger: logger, ProfileService: profileService}))
	r.GET("/", playgroundHandler())

	logger.InfoLog.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logger.ErrorLog.Fatal(r.Run())
}
