package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	mux := http.NewServeMux()

	execPath, _ := os.Executable()
	baseDir := filepath.Dir(execPath)
	staticDir := filepath.Join(baseDir, "static")

	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = "./docs/static"
	}

	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		staticDir = "./static"
	}

	fs := http.FileServer(http.Dir(staticDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}
		if r.URL.Path == "/docs" || r.URL.Path == "/docs/" {
			http.ServeFile(w, r, filepath.Join(staticDir, "docs.html"))
			return
		}
		http.NotFound(w, r)
	})

	log.Println("GoLaravel Documentation Server")
	log.Println("Serving on http://0.0.0.0:5000")
	log.Println("Static dir:", staticDir)

	if err := http.ListenAndServe("0.0.0.0:5000", mux); err != nil {
		log.Fatal(err)
	}
}
