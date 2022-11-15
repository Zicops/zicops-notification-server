package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	//"github.com/99designs/gqlgen/graphql/playground"
	"github.com/zicops/zicops-notification-server/graph"
	"github.com/zicops/zicops-notification-server/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
