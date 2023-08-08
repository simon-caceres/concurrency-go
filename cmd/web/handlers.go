package main

import (
	"final-project/data"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		if !match {
			msg := Message{
				To:      email,
				Subject: "Login attempt failed",
				Data:    "Invalid login credentials",
			}

			app.SendEmail(msg)
		}

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
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// TODO: - validate form data

	// create an user
	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}

	// insert the user into the database
	_, err = u.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Failed to create an account")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send an activation email
	url := fmt.Sprintf("http://localhost/activate?email=%s", u.Email)
	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       u.Email,
		Subject:  "Activate your account",
		Template: "Confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.SendEmail(msg)

	app.Session.Put(r.Context(), "flash", "Your account has been created. Please check your email to activate your account")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	// validate url parameters
	url := r.RequestURI
	testUrl := fmt.Sprintf("http://localhost%s", url)

	ok := VerifyToken(testUrl)

	if !ok {
		app.Session.Put(r.Context(), "error", "Invalid activation link")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// activate de account
	u, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u.Active = 1

	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to update user.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "Your account has been activated. Please login")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// generate an invoice

	// send an email with attachments

	// send an email with the invoice attached
}

func (app *Config) TestEmail(w http.ResponseWriter, r *http.Request) {
	m := Mail{
		Domain:      "localhost",
		Host:        "localhost",
		Port:        1025,
		Encryption:  "none",
		FromName:    "Test",
		FromAddress: "info@localhost",
		ErrorChan:   make(chan error),
	}

	msg := Message{
		To:      "me@localhost",
		Subject: "Test email",
		Data:    "This is a test email",
	}

	m.sendMail(msg, make(chan error))
}

func (app *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {
	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	// TODO: add a error page

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	// GET THE id of the plan that is chosen
	id := r.URL.Query().Get("id")

	planID, _ := strconv.Atoi(id)

	// get the plan from the database
	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to find a plan.")
		http.Redirect(w, r, "/menbers/plans", http.StatusSeeOther)
		return
	}

	// get the user from the session
	user, ok := app.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "Unable to find a user.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// generate an invoice
	app.Wait.Add(1)

	go func() {
		defer app.Wait.Done()

		// create an invoice
		invoice, err := app.GetInvoice(plan, user)
		if err != nil {
			app.ErrorChan <- err
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Invoice",
			Template: "invoice-email",
			Data:     invoice,
		}

		// send email
		app.SendEmail(msg)
	}()

	// send the email with the invoice attached

	// generate a manual

	// send the email with the manual attached

	// create a subscription

	// redirect
}

func (app *Config) GetInvoice(plan *data.Plan, u data.User) (string, error) {
	return plan.PlanAmountFormatted, nil
}
