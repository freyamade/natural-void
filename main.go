package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/context"
	// "github.com/gorilla/sessions"
	"github.com/go-chi/chi"

	"./go"
)

func main() {
	dao := naturalvoid.DAO{}
	err := dao.New()
	if err != nil {
		panic(err)
	}
	// Get the list of Stories from the DB
	stories := []naturalvoid.Story{}
	dao.DB.Find(&stories)
	data := struct {
		Stories []naturalvoid.Story
	}{
		Stories: stories,
	}
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := naturalvoid.ParseTemplatesInDir("templates")
		if err != nil {
			panic(err)
		}
		tmpl.ExecuteTemplate(w, "index.tmpl", &data)
	})
	// Serve the static files
	FileServer(r, "/static", http.Dir("./static"))
	http.ListenAndServe(":3333", context.ClearHandler(r))
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
