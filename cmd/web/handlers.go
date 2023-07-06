package main

import "net/http"

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
	// _ = app.Session.RenewToken(r.Context())

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

}

func (app *Config) RegistrationPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {

}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {

}
