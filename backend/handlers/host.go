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

func HostHandler(db *storage.Database ,w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "HX-Request, HX-Target, HX-Current-URL, Content-Type")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
	key, _ := db.CreateUniqueRoomId()
	hd := HostData{
		AccesKey:  key,
	}
	
    tmplPath := filepath.Join("templates", "fragment.html")

    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        log.Println("Error loading template:", err)
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }


    err = tmpl.ExecuteTemplate(w,"waitingRoom",hd)
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}

func Import_questions (db *storage.Database ,w http.ResponseWriter, r *http.Request) {

} 





















