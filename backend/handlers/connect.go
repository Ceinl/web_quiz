package handlers

import (
    "net/http"
	"fmt"
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
		fmt.Fprintf(w, "Hello from the Go backend!") // TODO: Send 4 buttons that return your answer
    }
