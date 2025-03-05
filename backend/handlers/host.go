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

func createRoom() (string,string){

	db, err := storage.CreaateDatabase("game.db")
	if err != nil {
		log.Fatal(err)
	}

	key, _ :=  db.CreateUniqueRoomId()
	return key, "question"
}


func HostHandler(w http.ResponseWriter, r *http.Request) {
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
        log.Println("Error loading template:", err)
        http.Error(w, "Error loading template", http.StatusInternalServerError)
        return
    }

	key,question := createRoom()

	data := PageData{
		Question: question,
		AccesKey: key, 
	}
	

    err = tmpl.ExecuteTemplate(w,"waitingRoom", data)
    if err != nil {
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}
























