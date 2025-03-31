package main

import (
	"flag"
	"fmt"

	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-migrate/migrate/v4"
	"github.com/kudras3r/CommentSystem/internal/api/graphql"
	"github.com/kudras3r/CommentSystem/internal/service"
	"github.com/kudras3r/CommentSystem/internal/storage"
	"github.com/kudras3r/CommentSystem/internal/storage/inmemory"
	"github.com/kudras3r/CommentSystem/internal/storage/migration"
	"github.com/kudras3r/CommentSystem/internal/storage/postgres"
	"github.com/kudras3r/CommentSystem/pkg/config"
	"github.com/kudras3r/CommentSystem/pkg/logger"
	"github.com/vektah/gqlparser/v2/ast"
)

const (
	DELAULTSTORAGE = "inmemory"
)

func main() {

	// cfg init
	config := config.Load()

	// logger init TODO
	log := logger.New(config.LogLevel)
	log.Infof("set level %s", log.Level)

	// storage init
	var storage storage.Storage

	storageKind := flag.String("storage", DELAULTSTORAGE, "storage kind: db / inmemory")
	flag.Parse()

	switch *storageKind {
	case "db":
		var err error
		pg, err := postgres.New(config.DB, log)
		if err != nil {
			log.Fatalf("cannot init db : %v", err)
		}
		defer pg.CloseConnection()

		// migrate
		if err := migration.MakeMigrations(pg.GetConnection()); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatalf("cannot migrate! : %v", err)
			}
		}
		log.Info("init db storage...")
		log.Info("migrate...")

		storage = pg

	case "inmemory":
		storage = inmemory.New(log)
		log.Info("init inmemory storage...")
	}

	log.Warnf("by default storage kind : %s, please choose it --storage=db / im if needed", DELAULTSTORAGE)
	if storage == nil {
		log.Fatal("storage is not initialized properly!")
	}

	// service init
	service := service.New(storage, log)
	resolver := &graphql.Resolver{
		Service: service,
	}
	log.Info("service up...")

	// srv generated
	srv := handler.New(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))

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

	// run
	log.Infof("connect to http://%s:%s/ for GraphQL playground", config.Server.Host, config.Server.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port), nil))
}
