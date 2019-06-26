package main

import (
	"github.com/freyamade/natural-void/go"
	"net/http"
)

func main() {
	r := naturalvoid.NewRouter()
	http.ListenAndServe(":3333", r)
}
