package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const defaultComposeTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Northbriton Publisher</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f4efe7;
      --panel: #fffdf8;
      --ink: #1f1b16;
      --muted: #65594b;
      --accent: #0b5c5a;
      --accent-2: #d1843b;
      --border: #d8c9b6;
      --danger: #b0452f;
      --shadow: 0 18px 40px rgba(44, 31, 13, 0.12);
      --radius: 18px;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      font-family: Georgia, "Iowan Old Style", serif;
      color: var(--ink);
      background:
        radial-gradient(circle at top left, rgba(209,132,59,0.18), transparent 28%),
        linear-gradient(180deg, #efe3d1 0%, var(--bg) 30%, #ede7de 100%);
    }
    main {
      max-width: 760px;
      margin: 0 auto;
      padding: 24px 16px 56px;
    }
    .panel {
      background: var(--panel);
      border: 1px solid var(--border);
      border-radius: var(--radius);
      box-shadow: var(--shadow);
      overflow: hidden;
    }
    .hero {
      padding: 22px 22px 18px;
      background:
        linear-gradient(135deg, rgba(11,92,90,0.96), rgba(17,85,122,0.88)),
        linear-gradient(45deg, rgba(255,255,255,0.14), transparent 60%);
      color: #fffaf2;
    }
    h1 {
      margin: 0 0 6px;
      font-size: clamp(2rem, 7vw, 3.2rem);
      line-height: 0.95;
      letter-spacing: -0.04em;
    }
    .hero p {
      margin: 0;
      max-width: 38rem;
      color: rgba(255,250,242,0.84);
    }
    .content {
      padding: 18px;
      display: grid;
      gap: 16px;
    }
    label {
      display: grid;
      gap: 8px;
      font-size: 0.95rem;
      color: var(--muted);
    }
    input, textarea, select {
      width: 100%;
      border: 1px solid var(--border);
      border-radius: 14px;
      font: inherit;
      color: var(--ink);
      background: #fff;
      padding: 14px;
    }
    textarea {
      min-height: 220px;
      resize: vertical;
    }
    .actions {
      display: flex;
      flex-wrap: wrap;
      gap: 12px;
      align-items: center;
      justify-content: space-between;
    }
    button, .button {
      border: 0;
      border-radius: 999px;
      padding: 14px 18px;
      font: inherit;
      text-decoration: none;
      cursor: pointer;
      background: var(--accent);
      color: white;
    }
    .secondary {
      background: rgba(11, 92, 90, 0.12);
      color: var(--accent);
    }
    .meta {
      display: flex;
      justify-content: space-between;
      gap: 12px;
      color: var(--muted);
      font-size: 0.95rem;
      align-items: center;
    }
    #status {
      min-height: 1.4rem;
      color: var(--muted);
    }
    #status.error {
      color: var(--danger);
    }
    #status.success {
      color: var(--accent);
    }
    .config {
      border-radius: 14px;
      background: #f8f0e6;
      border: 1px solid #e5cfb5;
      padding: 14px;
      color: #6f4b1c;
    }
    .login-panel {
      padding: 18px;
      display: grid;
      gap: 16px;
    }
    .small {
      font-size: 0.9rem;
      color: var(--muted);
    }
    @media (max-width: 640px) {
      .hero {
        padding: 18px 18px 16px;
      }
      .content {
        padding: 16px;
      }
      .actions {
        align-items: stretch;
      }
      .actions > * {
        width: 100%;
      }
      .meta {
        display: grid;
      }
    }
  </style>
</head>
<body>
  <main>
    <section class="panel">
      <div class="hero">
        <h1>Phone post, no laptop.</h1>
        <p>Write a short entry, attach photos, and let the service open a pull request or publish directly.</p>
      </div>
      {{ if .LoggedIn }}
      <form id="compose-form" class="content">
        {{ if .ConfigError }}
        <div class="config">{{ .ConfigError }}</div>
        {{ end }}
        <div class="meta">
          <span>Signed in as {{ .Email }}</span>
          <a class="button secondary" href="{{ .LogoutURL }}">Log out</a>
        </div>
        <input type="hidden" name="csrf_token" value="{{ .CSRFToken }}">
        <label>
          Title
          <input name="title" type="text" placeholder="Optional. Leave blank for an auto-generated title.">
        </label>
        <label>
          Entry
          <textarea name="body" placeholder="A quick note, a trail report, dinner photos, whatever you want to remember."></textarea>
        </label>
        <label>
          Tags
          <input name="tags" type="text" placeholder="family, photos, hiking">
        </label>
        <label>
          Photos
          <input name="images" type="file" accept="image/*" multiple>
        </label>
        <label>
          Publish mode
          <select name="publish_mode">
            <option value="pr">Create pull request</option>
            <option value="direct">Publish now</option>
          </select>
        </label>
        <div class="actions">
          <button type="submit" {{ if .ConfigError }}disabled{{ end }}>Send to blog</button>
          <div id="status" role="status" aria-live="polite"></div>
        </div>
      </form>
      {{ else }}
      <div class="login-panel">
        {{ if .ConfigError }}
        <div class="config">{{ .ConfigError }}</div>
        {{ else }}
        <p>Authentication is just one shared password. Images still upload directly from your phone to this service.</p>
        {{ if .LoginError }}
        <div class="config">{{ .LoginError }}</div>
        {{ end }}
        <form method="post" action="{{ .LoginURL }}" class="content">
          <label>
            Password
            <input name="password" type="password" autocomplete="current-password" placeholder="Configured by AUTH_PASSWORD">
          </label>
          <button type="submit">Sign in</button>
        </form>
        {{ end }}
        <div class="small">This is a single-author tool for {{ .AllowedEmail }}.</div>
      </div>
      {{ end }}
    </section>
  </main>
  {{ if .LoggedIn }}
  <script>
    const form = document.getElementById('compose-form');
    const status = document.getElementById('status');
    if (form) {
      form.addEventListener('submit', async (event) => {
        event.preventDefault();
        status.className = '';
        status.textContent = 'Submitting…';
        const formData = new FormData(form);
        try {
          const response = await fetch('{{ .PostURL }}', {
            method: 'POST',
            body: formData,
            credentials: 'same-origin'
          });
          const payload = await response.json();
          if (!response.ok) {
            throw new Error(payload.error || 'Request failed');
          }
          const link = payload.url ? '<a href="' + payload.url + '" target="_blank" rel="noreferrer">' + payload.url + '</a>' : '';
          status.className = 'success';
          status.innerHTML = payload.message + (link ? '<br>' + link : '');
          form.reset();
        } catch (error) {
          status.className = 'error';
          status.textContent = error.message;
        }
      });
    }
  </script>
  {{ end }}
</body>
</html>`

type Config struct {
	AllowedEmail          string
	BasePath              string
	CookieSecure          bool
	GitHubAuthorEmail     string
	GitHubAuthorName      string
	GitHubBaseBranch      string
	GitHubOwner           string
	GitHubRepo            string
	GitHubToken           string
	ImageTransformBaseURL string
	MaxUploadBytes        int64
	Port                  string
	PublicImageBaseURL    string
	AuthPassword          string
	SessionSecret         string
	TimeLocation          *time.Location
}

type server struct {
	cfg      Config
	github   *githubClient
	tmpl     *template.Template
	uploader uploader
}

type composePageData struct {
	AllowedEmail string
	CSRFToken    string
	ConfigError  string
	Email        string
	LoggedIn     bool
	LoginError   string
	LoginURL     string
	LogoutURL    string
	PostURL      string
}

type postResponse struct {
	Message string `json:"message"`
	URL     string `json:"url,omitempty"`
}

func loadConfig() Config {
	return Config{
		AllowedEmail:          strings.TrimSpace(os.Getenv("ALLOWED_EMAIL")),
		BasePath:              normalizeBasePath(os.Getenv("BASE_PATH")),
		CookieSecure:          parseBoolEnv("COOKIE_SECURE", true),
		GitHubAuthorEmail:     defaultString(os.Getenv("GITHUB_AUTHOR_EMAIL"), "paul.miller@gmail.com"),
		GitHubAuthorName:      defaultString(os.Getenv("GITHUB_AUTHOR_NAME"), "Paul Miller"),
		GitHubBaseBranch:      defaultString(os.Getenv("GITHUB_BASE_BRANCH"), "master"),
		GitHubOwner:           strings.TrimSpace(os.Getenv("GITHUB_OWNER")),
		GitHubRepo:            strings.TrimSpace(os.Getenv("GITHUB_REPO")),
		GitHubToken:           strings.TrimSpace(os.Getenv("GITHUB_TOKEN")),
		ImageTransformBaseURL: defaultString(os.Getenv("IMAGE_TRANSFORM_BASE_URL"), "https://images.northbriton.net/cdn-cgi/image/width=800"),
		MaxUploadBytes:        parseInt64Env("MAX_UPLOAD_BYTES", 64<<20),
		Port:                  defaultString(os.Getenv("PORT"), "8080"),
		PublicImageBaseURL:    strings.TrimRight(defaultString(os.Getenv("PUBLIC_IMAGE_BASE_URL"), "https://images.northbriton.net"), "/"),
		AuthPassword:          strings.TrimSpace(os.Getenv("AUTH_PASSWORD")),
		SessionSecret:         strings.TrimSpace(os.Getenv("SESSION_SECRET")),
		TimeLocation:          loadLocation(defaultString(os.Getenv("BLOG_TIMEZONE"), "America/Los_Angeles")),
	}
}

func newServer(cfg Config, u uploader) http.Handler {
	return &server{
		cfg:      cfg,
		github:   newGitHubClient(cfg),
		tmpl:     template.Must(template.New("compose").Parse(defaultComposeTemplate)),
		uploader: u,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routePath, ok := s.routePath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch routePath {
	case "/":
		s.handleRoot(w, r)
	case "/compose":
		s.handleCompose(w, r)
	case "/api/posts":
		s.handleCreatePost(w, r)
	case "/auth/login":
		s.handleLogin(w, r)
	case "/auth/logout":
		s.handleLogout(w, r)
	case "/mirror":
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		serveMirror(w, r, s.cfg, s.uploader)
	default:
		http.NotFound(w, r)
	}
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && strings.TrimSpace(r.FormValue("album")) != "" {
		serveMirror(w, r, s.cfg, s.uploader)
		return
	}
	if _, ok := s.currentSession(r); ok {
		http.Redirect(w, r, s.url("/compose"), http.StatusFound)
		return
	}
	s.renderComposePage(w, r, nil, "")
}

func (s *server) handleCompose(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET only", http.StatusMethodNotAllowed)
		return
	}

	session, ok := s.requireSession(w, r)
	if !ok {
		return
	}
	s.renderComposePage(w, r, session, "")
}

func (s *server) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	if !s.composeConfigured() {
		s.writeError(w, http.StatusServiceUnavailable, s.composeConfigError())
		return
	}

	session, ok := s.currentSession(r)
	if !ok {
		s.writeError(w, http.StatusUnauthorized, "sign in required")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, s.cfg.MaxUploadBytes)
	if err := r.ParseMultipartForm(s.cfg.MaxUploadBytes); err != nil {
		s.writeError(w, http.StatusBadRequest, "could not parse multipart form")
		return
	}
	if r.FormValue("csrf_token") != session.CSRFToken {
		s.writeError(w, http.StatusForbidden, "invalid CSRF token")
		return
	}

	req, err := s.buildPostRequest(r)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := s.publishPost(r.Context(), req)
	if err != nil {
		log.Printf("publish failed: %v", err)
		s.writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, postResponse{
		Message: result.Message,
		URL:     result.URL,
	})
}

func (s *server) buildPostRequest(r *http.Request) (postRequest, error) {
	req := postRequest{
		Title:       strings.TrimSpace(r.FormValue("title")),
		Body:        strings.TrimSpace(r.FormValue("body")),
		Tags:        parseTags(r.FormValue("tags")),
		PublishMode: parsePublishMode(r.FormValue("publish_mode")),
		Now:         time.Now().In(s.cfg.TimeLocation),
	}

	files := r.MultipartForm.File["images"]
	for _, header := range files {
		image, err := s.readImageUpload(header)
		if err != nil {
			return postRequest{}, err
		}
		req.Images = append(req.Images, image)
	}
	if err := req.Validate(); err != nil {
		return postRequest{}, err
	}
	return req, nil
}

func (s *server) publishPost(ctx context.Context, req postRequest) (publishResult, error) {
	rendered, err := buildPost(req)
	if err != nil {
		return publishResult{}, err
	}

	var imageURLs []string
	for _, image := range req.Images {
		fileName, err := newImageObjectName(req.Now, image.FileName, image.ContentType)
		if err != nil {
			return publishResult{}, err
		}
		if err := s.uploader.Put(ctx, fileName, image.ContentType, bytes.NewReader(image.Data)); err != nil {
			return publishResult{}, fmt.Errorf("upload %s: %w", image.FileName, err)
		}
		imageURLs = append(imageURLs, fmt.Sprintf("%s/%s", s.cfg.PublicImageBaseURL, fileName))
	}

	contentBytes := []byte(rendered.RenderMarkdown(imageURLs))
	postPath, err := s.github.EnsureUniquePath(ctx, rendered.FilePath(req.Now))
	if err != nil {
		return publishResult{}, fmt.Errorf("choose path: %w", err)
	}

	publishReq := githubPublishRequest{
		BaseBranch:  s.cfg.GitHubBaseBranch,
		Body:        contentBytes,
		CommitBody:  rendered.CommitMessage(req.PublishMode),
		FilePath:    postPath,
		PRBody:      rendered.PullRequestBody(),
		PRTitle:     rendered.PullRequestTitle(),
		PublishMode: req.PublishMode,
		Slug:        rendered.Slug,
	}
	return s.github.Publish(ctx, publishReq)
}

func (s *server) readImageUpload(header *multipart.FileHeader) (uploadedImage, error) {
	file, err := header.Open()
	if err != nil {
		return uploadedImage{}, fmt.Errorf("open upload %s: %w", header.Filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return uploadedImage{}, fmt.Errorf("read upload %s: %w", header.Filename, err)
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	if !strings.HasPrefix(contentType, "image/") {
		return uploadedImage{}, fmt.Errorf("%s is not an image", header.Filename)
	}
	return uploadedImage{
		ContentType: contentType,
		Data:        data,
		FileName:    header.Filename,
	}, nil
}

func (s *server) renderComposePage(w http.ResponseWriter, r *http.Request, session *sessionPayload, extraError string) {
	data := composePageData{
		AllowedEmail: s.authorIdentity(),
		ConfigError:  extraError,
		LoginError:   strings.TrimSpace(r.URL.Query().Get("error")),
		LoginURL:     s.url("/auth/login"),
		LogoutURL:    s.url("/auth/logout"),
		PostURL:      s.url("/api/posts"),
	}
	if data.ConfigError == "" {
		data.ConfigError = s.composeConfigError()
	}
	if session != nil {
		data.LoggedIn = true
		data.Email = session.Email
		data.CSRFToken = session.CSRFToken
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) routePath(requestPath string) (string, bool) {
	if requestPath == "" {
		requestPath = "/"
	}
	requestPath = path.Clean(requestPath)
	if s.cfg.BasePath == "" {
		return requestPath, true
	}
	if requestPath == s.cfg.BasePath {
		return "/", true
	}
	if !strings.HasPrefix(requestPath, s.cfg.BasePath+"/") {
		return "", false
	}
	trimmed := strings.TrimPrefix(requestPath, s.cfg.BasePath)
	if trimmed == "" {
		return "/", true
	}
	return trimmed, true
}

func (s *server) url(suffix string) string {
	if suffix == "" || suffix == "/" {
		if s.cfg.BasePath == "" {
			return "/"
		}
		return s.cfg.BasePath
	}
	if s.cfg.BasePath == "" {
		return suffix
	}
	return path.Clean(s.cfg.BasePath + "/" + strings.TrimPrefix(suffix, "/"))
}

func (s *server) composeConfigured() bool {
	return s.composeConfigError() == ""
}

func (s *server) composeConfigError() string {
	var missing []string
	if s.cfg.SessionSecret == "" {
		missing = append(missing, "SESSION_SECRET")
	}
	if s.cfg.AuthPassword == "" {
		missing = append(missing, "AUTH_PASSWORD")
	}
	if s.cfg.GitHubToken == "" || s.cfg.GitHubOwner == "" || s.cfg.GitHubRepo == "" {
		missing = append(missing, "GITHUB_TOKEN/GITHUB_OWNER/GITHUB_REPO")
	}
	if len(missing) == 0 {
		return ""
	}
	return "Publisher configuration is incomplete. Missing: " + strings.Join(missing, ", ")
}

func normalizeBasePath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "/" {
		return ""
	}
	value = "/" + strings.Trim(value, "/")
	return strings.TrimRight(value, "/")
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func parseBoolEnv(name string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseInt64Env(name string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func loadLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		return time.Local
	}
	return location
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (s *server) writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func requireMethod(r *http.Request, method string) error {
	if r.Method != method {
		return errors.New(method + " only")
	}
	return nil
}
