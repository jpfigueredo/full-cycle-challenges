package graphql

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func StartGraphQLServer(resolver *Resolver, port string) {
	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolver}))

	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))

	log.Printf("🚀 GraphQL server running on http://localhost:%s/", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ failed to start GraphQL server: %v", err)
	}
}
