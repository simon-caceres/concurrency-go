package main

import (
	"final-project/data"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var pageTests = []struct {
	name         string
	url          string
	statusCode   int
	handler      http.HandlerFunc
	sessionData  map[string]interface{}
	expectedHTML string
}{
	{
		name:         "home",
		url:          "/",
		statusCode:   200,
		handler:      testApp.HomePage,
		sessionData:  nil,
		expectedHTML: `<h1 class="mt-5">Home</h1>`,
	},
	{
		name:         "loginPage",
		url:          "/login",
		statusCode:   200,
		handler:      testApp.LoginPage,
		sessionData:  nil,
		expectedHTML: `<h1 class="mt-5">Login</h1>`,
	},
	{
		name:       "LogoutPage",
		url:        "/logout",
		statusCode: http.StatusSeeOther,
		handler:    testApp.LoginPage,
		sessionData: map[string]interface{}{
			"userID": 1,
			"user":   data.User{},
		},
		expectedHTML: `<h1 class="mt-5">Login</h1>`,
	},
}

func Test_pages(t *testing.T) {
	pathToTemplates = "./templates"

	for _, e := range pageTests {
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", e.url, nil)

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		if len(e.sessionData) > 0 {
			for key, value := range e.sessionData {
				testApp.Session.Put(ctx, key, value)
			}
		}

		e.handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("want %d; got %d", http.StatusOK, rr.Code)
		}

		if len(e.expectedHTML) > 0 {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("want %s; got %s", e.expectedHTML, html)
			}
		}
	}
}

func TestConfig_Home(t *testing.T) {
	pathToTemplates = "./templates"

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	handdler := http.HandlerFunc(testApp.HomePage)

	handdler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rr.Code)
	}
}

func TestConfig_PostLoginPage(t *testing.T) {
	pathToTemplates = "./templates"

	postedData := url.Values{
		"email": {
			"admin@example.com",
		},
		"password": {
			"password",
		},
	}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(testApp.PostLoginPage)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("want %d; got %d", http.StatusSeeOther, rr.Code)
	}

	if !testApp.Session.Exists(ctx, "userID") {
		t.Error("session does not exist")
	}
}
