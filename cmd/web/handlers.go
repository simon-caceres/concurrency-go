package main

import (
	"fmt"
	"net/http"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) AboutPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())

	// parse form data
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get the email and pasword from the form
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// tried to get the user from the database
	user, err := app.Models.User.GetByEmail(email)

	// if the user is not found
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// check if the password is correct
	match, err := user.PasswordMatches(password)

	// if the password is not correct
	if err != nil || !match {
		fmt.Println("password is not correct")
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// if the password is correct
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	app.Session.Put(r.Context(), "flash", "Login successful")

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (app *Config) LogoutPage(w http.ResponseWriter, r *http.Request) {
	// clean up session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())
	// redirect to home page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegistrationPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	// create an user

	// send an activation email

	// subscribe the user to an account
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url parameters

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached
}
