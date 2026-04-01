package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	sessionCookieName = "album_session"
	sessionLifetime   = 14 * 24 * time.Hour
)

type sessionPayload struct {
	CSRFToken string `json:"csrf_token"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := requireMethod(r, http.MethodPost); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	if s.cfg.SessionSecret == "" || s.cfg.AuthPassword == "" {
		http.Error(w, s.composeConfigError(), http.StatusServiceUnavailable)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "could not parse login form", http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if subtle.ConstantTimeCompare([]byte(password), []byte(s.cfg.AuthPassword)) != 1 {
		http.Redirect(w, r, s.url("/")+"?error=wrong+password", http.StatusFound)
		return
	}

	csrfToken, err := randomToken(32)
	if err != nil {
		http.Error(w, "could not create session", http.StatusInternalServerError)
		return
	}
	session := sessionPayload{
		CSRFToken: csrfToken,
		Email:     s.authorIdentity(),
		ExpiresAt: time.Now().Add(sessionLifetime).Unix(),
	}
	if err := s.setSignedCookie(w, sessionCookieName, session, sessionLifetime); err != nil {
		http.Error(w, "could not persist session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, s.url("/compose"), http.StatusFound)
}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if err := requireMethod(r, http.MethodGet); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}
	s.clearCookie(w, sessionCookieName)
	http.Redirect(w, r, s.url("/"), http.StatusFound)
}

func (s *server) currentSession(r *http.Request) (*sessionPayload, bool) {
	var payload sessionPayload
	if err := s.readSignedCookie(r, sessionCookieName, &payload); err != nil {
		return nil, false
	}
	if time.Now().Unix() > payload.ExpiresAt {
		return nil, false
	}
	if payload.Email != s.authorIdentity() {
		return nil, false
	}
	return &payload, true
}

func (s *server) requireSession(w http.ResponseWriter, r *http.Request) (*sessionPayload, bool) {
	session, ok := s.currentSession(r)
	if ok {
		return session, true
	}
	http.Redirect(w, r, s.url("/"), http.StatusFound)
	return nil, false
}

func (s *server) authorIdentity() string {
	for _, value := range []string{s.cfg.AllowedEmail, s.cfg.GitHubAuthorEmail, s.cfg.GitHubAuthorName, "blog-author"} {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return "blog-author"
}

func (s *server) setSignedCookie(w http.ResponseWriter, name string, payload sessionPayload, maxAge time.Duration) error {
	encoded, err := s.signPayload(payload)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    encoded,
		HttpOnly: true,
		MaxAge:   int(maxAge.Seconds()),
		Path:     cookiePath(s.cfg.BasePath),
		SameSite: http.SameSiteLaxMode,
		Secure:   s.cfg.CookieSecure,
	})
	return nil
}

func (s *server) clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     cookiePath(s.cfg.BasePath),
		SameSite: http.SameSiteLaxMode,
		Secure:   s.cfg.CookieSecure,
	})
}

func (s *server) readSignedCookie(r *http.Request, name string, out *sessionPayload) error {
	if s.cfg.SessionSecret == "" {
		return errors.New("session secret not configured")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return err
	}
	payloadBytes, err := decodeSignedValue(cookie.Value, []byte(s.cfg.SessionSecret))
	if err != nil {
		return err
	}
	return json.Unmarshal(payloadBytes, out)
}

func (s *server) signPayload(payload sessionPayload) (string, error) {
	if s.cfg.SessionSecret == "" {
		return "", errors.New("session secret not configured")
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return encodeSignedValue(raw, []byte(s.cfg.SessionSecret)), nil
}

func randomToken(bytesLen int) (string, error) {
	raw := make([]byte, bytesLen)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func encodeSignedValue(payload []byte, secret []byte) string {
	sig := computeHMAC(payload, secret)
	return base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func decodeSignedValue(value string, secret []byte) ([]byte, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return nil, errors.New("bad cookie format")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	expected := computeHMAC(payload, secret)
	if !hmac.Equal(signature, expected) {
		return nil, errors.New("bad cookie signature")
	}
	return payload, nil
}

func computeHMAC(payload []byte, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(payload)
	return mac.Sum(nil)
}

func cookiePath(basePath string) string {
	if basePath == "" {
		return "/"
	}
	return basePath
}
