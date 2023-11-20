package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gin-gonic/gin"
	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/graph"
	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Environment struct {
	DatabaseUrl             string `json:"DATABASE_URL"`
	Auth0Domain             string `json:"AUTH0_DOMAIN"`
	Auth0ClientId           string `json:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret       string `json:"AUTH0_CLIENT_SECRET"`
	Auth0AudienceApi        string `json:"AUTH0_AUDIENCE_API"`
	SendgridApiKey          string `json:"SENDGRID_API_KEY"`
	SendgridListIdAccount   string `json:"SENDGRID_LIST_ID_ACCOUNT"`
	SendgridListIdMarketing string `json:"SENDGRID_LIST_ID_MARKETING"`
}

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

	// Load environment variables
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.ErrorLog.Fatalln(err)
	}

	svc := secretsmanager.NewFromConfig(cfg)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("AWS_SECRET_NAME")),
	}

	result, err := svc.GetSecretValue(context.Background(), input)
	if err != nil {
		logger.ErrorLog.Fatalln(err)
	}

	environment := &Environment{}
	if err := json.Unmarshal([]byte(*result.SecretString), environment); err != nil {
		logger.ErrorLog.Fatalln(err)
	}

	// Initialize services
	db, err := sql.Open("postgres", environment.DatabaseUrl)
	if err != nil {
		logger.ErrorLog.Fatalln(err)
	}

	auth0Config := &auth.Auth0Config{
		Auth0Domain:       environment.Auth0Domain,
		Auth0ClientId:     environment.Auth0ClientId,
		Auth0ClientSecret: environment.Auth0ClientSecret,
		Auth0AudienceApi:  environment.Auth0AudienceApi,
	}

	mailList := misc.NewMailList(environment.SendgridApiKey, environment.SendgridListIdAccount, environment.SendgridListIdMarketing)

	profileService := profile.NewProfileService(db, mailList)

	// Initialize handlers
	r := gin.Default()
	r.Use(misc.GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler(logger, auth0Config, &graph.Resolver{Logger: logger, ProfileService: profileService}))
	r.GET("/", playgroundHandler())

	logger.InfoLog.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logger.ErrorLog.Fatal(r.Run())
}
