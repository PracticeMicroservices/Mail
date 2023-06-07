package main

import (
	"log"
	"mail/cmd/api/controllers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type App struct {
	Mail controllers.Mail
}

func NewApp() *App {
	return &App{
		Mail: controllers.NewMailController(),
	}
}

func (a *App) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	mux.Use(middleware.Heartbeat("/healthCheck"))

	mux.Post("/send", a.Mail.SendMail)
	return mux
}

func (a *App) serve() {
	//define server
	srv := &http.Server{
		Addr:    ":80",
		Handler: a.routes(),
	}

	//start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
