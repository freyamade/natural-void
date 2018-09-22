package main

import (
	"./go"
	"net/http"
)

func main() {
	r := naturalvoid.NewRouter()
	http.ListenAndServe(":3333", r)
}
