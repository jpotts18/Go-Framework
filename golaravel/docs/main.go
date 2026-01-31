package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golaravel/app"
	"golaravel/app/http/middleware"
	"golaravel/app/http/request"
	"golaravel/app/http/response"
)

func main() {
	application := app.New(".")

	staticDir := findStaticDir()
	log.Println("Static directory:", staticDir)

	application.Router.Use(middleware.Logger())
	application.Router.Use(middleware.Recovery())
	application.Router.Use(middleware.SecureHeaders())

	application.Router.Get("/", func(req *request.Request, res *response.Response) error {
		content, err := os.ReadFile(filepath.Join(staticDir, "index.html"))
		if err != nil {
			return res.NotFound("Page not found")
		}
		res.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		return res.HTML(string(content))
	})

	application.Router.Get("/docs", func(req *request.Request, res *response.Response) error {
		content, err := os.ReadFile(filepath.Join(staticDir, "docs.html"))
		if err != nil {
			return res.NotFound("Page not found")
		}
		res.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		return res.HTML(string(content))
	})

	mux := http.NewServeMux()

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimPrefix(r.URL.Path, "/static/")
		fullPath := filepath.Join(staticDir, filePath)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		contentType := getContentType(filePath)
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Cache-Control", "public, max-age=3600")

		file, err := os.Open(fullPath)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		io.Copy(w, file)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		application.Router.ServeHTTP(w, r)
	})

	log.Println("GoLaravel Documentation Server (powered by GoLaravel)")
	log.Println("Serving on http://0.0.0.0:5000")

	if err := http.ListenAndServe("0.0.0.0:5000", mux); err != nil {
		log.Fatal(err)
	}
}

func findStaticDir() string {
	candidates := []string{
		"./golaravel/docs/static",
		"./docs/static",
		"./static",
	}

	execPath, _ := os.Executable()
	baseDir := filepath.Dir(execPath)
	candidates = append([]string{filepath.Join(baseDir, "static")}, candidates...)

	for _, dir := range candidates {
		if _, err := os.Stat(dir); err == nil {
			return dir
		}
	}

	return "./static"
}

func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".html":
		return "text/html; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	default:
		return "application/octet-stream"
	}
}
