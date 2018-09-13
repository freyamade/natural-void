package main

import (
	"./go"
	"github.com/gorilla/context"
	"net/http"
)

func main() {
	r := naturalvoid.NewRouter()
	http.ListenAndServe(":3333", context.ClearHandler(r))
}
