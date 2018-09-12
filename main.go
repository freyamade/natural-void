package main

import (
	"net/http"
	"github.com/gorilla/context"
	"./go"
)

func main() {
	r := naturalvoid.NewRouter()
	http.ListenAndServe(":3333", context.ClearHandler(r))
}


