package main

import (
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/glimpzio/backend/auth"
	"github.com/glimpzio/backend/connections"
	"github.com/glimpzio/backend/graph"
	"github.com/glimpzio/backend/misc"
	"github.com/glimpzio/backend/profile"
	"github.com/joho/godotenv"
)

type environment struct {
	Auth0Domain             string `json:"AUTH0_DOMAIN"`
	Auth0ClientId           string `json:"AUTH0_CLIENT_ID"`
	Auth0AudienceApi        string `json:"AUTH0_AUDIENCE_API"`
	SendgridApiKey          string `json:"SENDGRID_API_KEY"`
	SendgridListIdAccount   string `json:"SENDGRID_LIST_ID_ACCOUNT"`
	SendgridListIdMarketing string `json:"SENDGRID_LIST_ID_MARKETING"`
	SendgridSenderName      string `json:"SENDGRID_SENDER_NAME"`
	SendgridSenderEmail     string `json:"SENDGRID_SENDER_EMAIL"`
	LandingBaseUrl          string `json:"LANDING_BASE_URL"`
	AppBaseUrl              string `json:"APP_BASE_URL"`
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
	env := &environment{}
	if err := misc.LoadSecret(os.Getenv("AWS_SECRET_NAME"), env); err != nil {
		logger.ErrorLog.Fatalln(err)
	}

	// Initialize services
	db, err := misc.LoadDatabaseFromSecret(os.Getenv("DB_SECRET_NAME"))
	if err != nil {
		logger.ErrorLog.Fatalln(err)
	}
	defer db.Close()

	auth0Config := &auth.Auth0Config{
		Auth0Domain:      env.Auth0Domain,
		Auth0ClientId:    env.Auth0ClientId,
		Auth0AudienceApi: env.Auth0AudienceApi,
	}

	mailList := misc.NewMailList(env.SendgridApiKey, env.SendgridSenderName, env.SendgridSenderEmail, env.SendgridListIdAccount, env.SendgridListIdMarketing)

	profileService := profile.NewProfileService(db, db, mailList)
	connectionService := connections.NewConnectionService(db, db, mailList, profileService, env.LandingBaseUrl)

	// Initialize handlers
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{env.LandingBaseUrl, env.AppBaseUrl},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	r.Use(misc.GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler(logger, auth0Config, &graph.Resolver{Logger: logger, ProfileService: profileService, ConnectionService: connectionService, Auth0Config: auth0Config}))
	r.GET("/", playgroundHandler())

	logger.InfoLog.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logger.ErrorLog.Fatal(r.Run())
}
