package handlers

import (
    "html/template"
    "log"
    "net/http"
    "path/filepath"
	"quiz/storage"
)


type PageData struct{
	Question string
	AccesKey string
}

var GameInstance PageData




func Game (db *storage.Database ,w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "HX-Request, HX-Target, HX-Current-URL, Content-Type")

    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

	GameInstance.AccesKey = r.URL.Query().Get("acces_key")

    tmplPath := filepath.Join("templates", "fragment.html")

    tmpl, err := template.ParseFiles(tmplPath)
    if err != nil {
        log.Println("Error loading template:", err)
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
	}
	

	err = tmpl.ExecuteTemplate(w,"Game",GameInstance)

}	
