package main

import "net/http"

func (app *Config) SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func (app *Config) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "userID") {
			app.Session.Put(r.Context(), "error", "Access denied")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
