package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"quiz/storage"
)

type HostData struct {
	AccesKey string
}

func HostHandler(db *storage.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "HX-Request, HX-Target, HX-Current-URL, Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	key, _ := db.CreateUniqueRoomId()
	hd := HostData{
		AccesKey: key,
	}

	tmplPath := filepath.Join("templates", "fragment.html")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println("Error loading template:", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "waitingRoom", hd)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func ImportQuestionsHandler(db *storage.Database, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "HX-Request, HX-Target, HX-Current-URL, Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check content type before parsing form
	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)

	if err := db.ProcessFileUpload(contentType, r.ContentLength); err != nil {
		log.Printf("File validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Now parse the multipart form
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "File size error", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Printf("File error: %v", err)
		http.Error(w, "File error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	log.Printf("Processing file: %s, size: %d bytes", fileHeader.Filename, fileHeader.Size)

	if _, err := storage.Reader(file, db); err != nil {
		log.Printf("Error processing file: %v", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File imported successfully"))
}
