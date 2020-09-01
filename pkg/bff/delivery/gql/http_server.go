package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/napnap/middleware"
	internalMiddleware "github.com/jasonsoft/starter/internal/pkg/middleware"
	eventProto "github.com/jasonsoft/starter/pkg/event/proto"
	walletProto "github.com/jasonsoft/starter/pkg/wallet/proto"
	temporalClient "go.temporal.io/sdk/client"
)

func verifyOrigin(origin string) bool {
	return true
}

func NewHTTPServer(eventClient eventProto.EventServiceClient, walletClient walletProto.WalletServiceClient, temporalClient temporalClient.Client) *napnap.NapNap {
	resolver := Resolver{
		eventClient:    eventClient,
		walletClient:   walletClient,
		temporalClient: temporalClient,
	}

	nap := napnap.New()
	nap.Use(internalMiddleware.NewRequestIDMW())
	nap.Use(internalMiddleware.NewTracerMW())
	nap.Use(internalMiddleware.NewLoggerMW())

	// turn on CORS feature
	options := middleware.Options{}
	options.AllowOriginFunc = verifyOrigin
	options.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTION", "HEAD"}
	options.AllowedHeaders = []string{"*", "Authorization", "Content-Type", "Origin", "Content-Length", "accept"}
	nap.Use(middleware.NewCors(options))

	nap.Get("/", func(c *napnap.Context) error {
		return c.String(200, "Hello World")
	})

	nap.Get("/ping", func(c *napnap.Context) error {
		return c.String(http.StatusOK, "pong!!!")
	})

	nap.Get("/playground", napnap.WrapHandler(handler.Playground("GraphQL playground", "/graphql")))

	nap.Post("/graphql", napnap.WrapHandler(handler.GraphQL(
		NewExecutableSchema(Config{Resolvers: &resolver}),
		CentralizedError(),
		RecoverError(),
	)))

	return nap
}
