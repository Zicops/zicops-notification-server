package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph"
	"github.com/zicops/zicops-notification-server/graph/generated"

	"github.com/zicops/zicops-notification-server/jwt"
)

const defaultPort = "8094"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/healthz"))
	router.Use(Middleware())
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":"+port, router))
	defer global.Client.Close()
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := jwt.GetToken(r)
			if token == "" {
				http.Error(w, "Unauthorized: Bad request or authorization details, invalid token", http.StatusUnauthorized)
				return
			}
			// put it in context
			ctx := context.WithValue(r.Context(), "token", token)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
