package main

import (
    "net/http"
    "github.com/go-chi/chi"
)

func main() {
    http.Handle("/", http.FileServer(http.Dir("./dist")))
    http.ListenAndServe(":3000", nil)
}
