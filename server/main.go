package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"graphql-poc/internal/db"
	"graphql-poc/internal/graph"
	"graphql-poc/internal/httpserver"
	"graphql-poc/internal/todo"
)

const defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/graphql_poc?sslmode=disable"

func main() {
	ctx := context.Background()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = defaultDatabaseURL
	}

	pool, err := db.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatal(err)
	}

	schema, err := graph.NewSchema(todo.NewStore(pool))
	if err != nil {
		log.Fatalf("failed to build schema: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("GraphQL server running at http://localhost:%s/graphql", port)
	log.Fatal(http.ListenAndServe(":"+port, httpserver.New(&schema)))
}
