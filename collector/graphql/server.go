package graphql

import (
	"context"
	"log"
	"net/http"
	"os"
	"oversee/collector/audit"
	"oversee/collector/graphql/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

type GraphqlAPIServer struct {
	searchService *audit.SearchService
	server        *http.Server
}

func NewGraphqlAPIServer(searchService *audit.SearchService) *GraphqlAPIServer {
	return &GraphqlAPIServer{
		searchService: searchService,
	}
}

func (g *GraphqlAPIServer) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		SearchService: g.searchService,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	g.server = &http.Server{Addr: ":" + port}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	return g.server.ListenAndServe()
}

func (g *GraphqlAPIServer) Shutdown(ctx context.Context) error {
	if g.server != nil {
		return g.server.Shutdown(ctx)
	}
	return nil
}
