package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	graph "github.com/kudras3r/CommentSystem/internal/api/graphql"
	"github.com/kudras3r/CommentSystem/internal/storage/migrate"
	"github.com/kudras3r/CommentSystem/internal/storage/postgres"
	"github.com/kudras3r/CommentSystem/pkg/config"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {

	// cfg init
	config := config.Load()

	// logger init TODO

	// storage init
	storage, err := postgres.New(config.DB)
	if err != nil {

	}
	defer storage.CloseConnection()

	// migrations
	if err := migrate.CreateTables(storage.GetConnection()); err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Println("migrate...")

	// srv generated
	resolver := &graph.Resolver{
		Storage: storage,
	}
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

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

	log.Printf("connect to http://%s:%s/ for GraphQL playground", config.Server.Host, config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port), nil))
}
