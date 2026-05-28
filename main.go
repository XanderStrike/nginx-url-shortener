package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	port      = 9090
	configDir = "/etc/nginx/redirects"
	baseURL   = "http://localhost:8080"
)

var idLength = 5

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	if s := os.Getenv("ID_LENGTH"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			idLength = n
		}
	}
}

func generateID() string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func sanitizeURL(u string) string {
	u = strings.ReplaceAll(u, "\n", "")
	u = strings.ReplaceAll(u, "\r", "")
	u = strings.ReplaceAll(u, "'", "")
	return u
}

func writeRedirectConfig(id, url string) error {
	safe := sanitizeURL(url)
	content := fmt.Sprintf("location = /%s {\n    return 301 '%s';\n}\n", id, safe)
	return os.WriteFile(filepath.Join(configDir, id+".conf"), []byte(content), 0644)
}

func reloadNginx() error {
	cmd := exec.Command("nginx", "-s", "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("reload failed: %s: %w", string(output), err)
	}
	return nil
}

var pageTmpl = template.Must(template.New("page").Parse(`<!DOCTYPE html>
<html>
<head>
<title>URL Shortener</title>
<style>
body { font-family: sans-serif; max-width: 600px; margin: 40px auto; padding: 0 20px; }
input[type=text] { width: 100%; padding: 8px; box-sizing: border-box; margin: 8px 0; }
button { padding: 8px 20px; cursor: pointer; }
.result { margin-top: 20px; padding: 16px; background: #f0f0f0; border-radius: 4px; }
.result a { word-break: break-all; }
</style>
</head>
<body>
<h1>URL Shortener</h1>
<form method="POST">
	<label>Paste your URL:</label><br>
	<input type="text" name="url" placeholder="https://example.com/very/long/url" required>
	<button type="submit">Shorten</button>
</form>
{{- if .ShortURL }}
<div class="result">
	<p>Short URL: <a href="{{.ShortURL}}">{{.ShortURL}}</a></p>
	<p>Redirects to: {{.OriginalURL}}</p>
</div>
{{- end }}
{{- if .Error }}
<div class="result" style="background:#ffe0e0">
	<p>Error: {{.Error}}</p>
</div>
{{- end }}
</body>
</html>`))

type pageData struct {
	ShortURL    string
	OriginalURL string
	Error       string
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := pageData{}

	if r.Method == http.MethodPost {
		url := r.FormValue("url")
		if url == "" {
			data.Error = "URL is required"
		} else {
			id := generateID()
			if err := writeRedirectConfig(id, url); err != nil {
				data.Error = fmt.Sprintf("Failed to write config: %v", err)
			} else if err := reloadNginx(); err != nil {
				data.Error = fmt.Sprintf("Config written but reload failed: %v", err)
			} else {
				data.ShortURL = baseURL + "/" + id
				data.OriginalURL = url
			}
		}
	}

	w.Header().Set("Content-Type", "text/html")
	pageTmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("listening on :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
