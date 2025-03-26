package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"quiz/storage"
)

func Connect(db *storage.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "HX-Request, HX-Target, HX-Current-URL, Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	tmplPath := filepath.Join("templates", "fragment.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error loading template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a data structure for the template
	data := struct {
		Title string
	}{
		Title: "Connect to Game",
	}

	if err := tmpl.ExecuteTemplate(w, "connect", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
