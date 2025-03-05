package main

import (
	"fmt"
	"net/http"
	"quiz/handlers"

)


type Room struct {
	AccesKey string
	Users []string
	Questions []questions
}

type questions struct {
	Question string
}

func createRoom () Room {
	q:= questions {"test"}
	qq := []questions{q}
	u := []string {"Dima","Vlad"}

	return Room{"qwerty", u, qq} //TODO: Make this a sqlite db and way to create/delete from db

}



func main() {

//	room := createRoom()
	fmt.Println("log")

	http.HandleFunc("/api/host", handlers.HostHandler)

    http.HandleFunc("/api/connect", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") 
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") 
        w.Header().Set("Access-Control-Allow-Headers", "HX-Request, HX-Target, HX-Current-URL, Content-Type") 
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
		fmt.Fprintf(w, "Hello from the Go backend!") // TODO: Send 4 buttons that return your answer
    })



    fmt.Println("Backend running at http://localhost:8081")
    http.ListenAndServe(":8081", nil)
} 
