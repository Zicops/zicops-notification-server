package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/robfig/cron/v3"
	"github.com/zicops/zicops-notification-server/global"
	"github.com/zicops/zicops-notification-server/graph"
	"github.com/zicops/zicops-notification-server/graph/generated"
	"github.com/zicops/zicops-notification-server/handlers"

	"github.com/zicops/zicops-notification-server/jwt"
)

const defaultPort = "8094"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	os.Setenv("SENDGRID_API_KEY", "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")
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/healthz"))
	router.Use(Middleware())
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":"+port, router))
	defer global.Client.Close()

	c := cron.New()
	c.AddFunc("0 2 * * 5", sch)

	c.Start()
}

func sch() {
	_, err := handlers.Auth_tokens(global.Ct)
	if err != nil {
		log.Printf("Got error while deleting tokens: %v", err)
	}
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

			//get fcm token
			fcm := r.Header.Get("fcm-token")
			if fcm == "" {
				http.Error(w, "Please mention fcm-token", http.StatusUnauthorized)
			}
			//put that in context
			ctxFcm := context.WithValue(ctx, "fcm-token", fcm)

			//accept tenant from headers
			ten := r.Header.Get("tenant")
			if ten == "" {
				http.Error(w, "Please mention tenant ", http.StatusUnauthorized)
			}
			//put that in context
			ctxTenant := context.WithValue(ctxFcm, "tenant", ten)

			// and call the next with our new context
			r = r.WithContext(ctxTenant)
			next.ServeHTTP(w, r)
		})
	}
}
