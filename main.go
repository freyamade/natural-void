package main

import (
	"./go"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"strings"
)

func main() {
	// Define a test StoryData slice
	stories := []naturalvoid.Story{{
		Name:      "A Simple Trip to Waterdeep: Dread",
		ShortName: "ASTTW: Dread",
		Slug:      "asttw-dread",
		Description: []string{
			"Chapter 1 of our 5 chapter epic which follows our heroes Bran, Gundham, Lyra, Jake, and Viper on their respective trips to Waterdeep.",
			"During one night of particularly heavy fog, these five people and their coach driver get taken to the land of Barovia, ruled by Strahd von Zarovich.",
			"Will this group of 5 strangers be able to band together, overcome this situation and rescue their coach driver? Only time will tell.",
		},
	}, {
		Name:      "A Simple Trip to Waterdeep: Reminiscence",
		ShortName: "ASTTW: Reminiscence",
		Slug:      "asttw-reminiscence",
		Description: []string{
			"Chapter 2 of our 5 chapter epic which follows our heroes Bran, Gundham, Lyra, Jake, and Viper on their respective trips to Waterdeep.",
			"After freeing themselves from the fog and getting back on the road, more weird scenarios begin to unfold.",
		},
	}}
	data := struct{
		Stories []naturalvoid.Story
	}{
        Stories: stories,
    }
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.tmpl")
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, data)
	})
	// Serve the static files
	FileServer(r, "/static", http.Dir("./static"))
	http.ListenAndServe(":3333", r)
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
