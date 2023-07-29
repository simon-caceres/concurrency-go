package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *Config) routes() http.Handler {
	// create a router

	mux := chi.NewRouter()

	// set up middleware

	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	// define application routes

	mux.Get("/", app.HomePage)

	mux.Get("/about", app.AboutPage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.LogoutPage)
	mux.Get("/register", app.RegistrationPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate-account", app.ActivateAccount)

	mux.Get("/test-email", app.TestEmail)

	return mux
}
