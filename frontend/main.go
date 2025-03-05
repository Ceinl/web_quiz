package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Роздаємо статичні файли з директорії
    http.Handle("/", http.FileServer(http.Dir(".")))

    fmt.Println("Фронтенд запущено на http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
