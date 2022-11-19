package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/zicops/zicops-notification-server/graph"
	"github.com/zicops/zicops-notification-server/graph/generated"
)

const defaultPort = "8094"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/query", srv)
	http.HandleFunc("/healthz", HealthCheckHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
