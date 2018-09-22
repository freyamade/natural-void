package main

import (
	"net/http"
    "./go"
)

func main() {
	r := naturalvoid.NewRouter()
	http.ListenAndServe(":3333", r)
}
