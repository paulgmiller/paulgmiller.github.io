package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandleLoginSetsSessionOnCorrectPassword(t *testing.T) {
	s := &server{
		cfg: Config{
			AllowedEmail:  "paul@example.com",
			AuthPassword:  "secret",
			CookieSecure:  false,
			SessionSecret: "session-secret",
		},
	}

	form := url.Values{}
	form.Set("password", "secret")
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	s.handleLogin(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if got, want := res.StatusCode, http.StatusFound; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if location := res.Header.Get("Location"); location != "/compose" {
		t.Fatalf("location = %q, want /compose", location)
	}
	if len(res.Cookies()) == 0 || res.Cookies()[0].Name != sessionCookieName {
		t.Fatalf("expected session cookie, got %#v", res.Cookies())
	}
}

func TestHandleLoginRejectsWrongPassword(t *testing.T) {
	s := &server{
		cfg: Config{
			AuthPassword:  "secret",
			CookieSecure:  false,
			SessionSecret: "session-secret",
		},
	}

	form := url.Values{}
	form.Set("password", "wrong")
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	s.handleLogin(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if got, want := res.StatusCode, http.StatusFound; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if location := res.Header.Get("Location"); location != "/?error=wrong+password" {
		t.Fatalf("location = %q, want /?error=wrong+password", location)
	}
}
