package gql

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	internalMiddleware "github.com/nite-coder/blackbear-demo/internal/pkg/middleware"
	eventProto "github.com/nite-coder/blackbear-demo/pkg/event/proto"
	walletProto "github.com/nite-coder/blackbear-demo/pkg/wallet/proto"
	"github.com/nite-coder/blackbear/pkg/web"
	"github.com/nite-coder/blackbear/pkg/web/middleware"
	temporalClient "go.temporal.io/sdk/client"
)

func verifyOrigin(origin string) bool {
	return true
}

func NewHTTPServer(eventClient eventProto.EventServiceClient, walletClient walletProto.WalletServiceClient, temporalClient temporalClient.Client) *web.WebServer {
	resolver := Resolver{
		eventClient:    eventClient,
		walletClient:   walletClient,
		temporalClient: temporalClient,
	}

	s := web.NewServer()
	s.Use(internalMiddleware.NewRequestIDMW())
	s.Use(internalMiddleware.NewTracerMW())
	s.Use(internalMiddleware.NewLoggerMW())

	// turn on CORS feature
	options := middleware.Options{}
	options.AllowOriginFunc = verifyOrigin
	options.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTION", "HEAD"}
	options.AllowedHeaders = []string{"*", "Authorization", "Content-Type", "Origin", "Content-Length", "accept"}
	s.Use(middleware.NewCors(options))

	s.Get("/", func(c *web.Context) error {
		return c.String(200, "Hello World")
	})

	s.Get("/ping", func(c *web.Context) error {
		return c.String(http.StatusOK, "pong!!!")
	})

	s.Get("/playground", web.WrapHandler(handler.Playground("GraphQL playground", "/graphql")))

	s.Post("/graphql", web.WrapHandler(handler.GraphQL(
		NewExecutableSchema(Config{Resolvers: &resolver}),
		CentralizedError(),
		RecoverError(),
	)))

	return s
}
