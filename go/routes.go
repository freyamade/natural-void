package naturalvoid

import (
    "net/http"
    "strings"
    "github.com/go-chi/chi"
    // "github.com/gorilla/sessions"
)

func NewRouter() (chi.Router) {
    r := chi.NewRouter()
    r.Get("/", Index)
    // Serve the static files
    fileServer(r, "/static", http.Dir("./static"))
    return r
}

// Define Route Handlers here
func Index(w http.ResponseWriter, r *http.Request) {
    // Get the list of Stories from the DB
    dao := GetDAO()
    stories := []Story{}
    dao.DB.Find(&stories)
    data := struct {
        Stories []Story
    }{
        Stories: stories,
    }
    // Generate and parse the templates
    tmpl, err := ParseTemplatesInDir("templates")
    if err != nil {
        panic(err)
    }
    tmpl.ExecuteTemplate(w, "index.tmpl", &data)
}

// Handle logging in of a user by checking against LDAP
func Login(w http.ResponseWriter, r *http.Request) {}

// Handle logging out of a logged in user
func Logout(w http.ResponseWriter, r *http.Request) {}

// Show the list of episodes in order of newest first
func Episodes(w http.ResponseWriter, r *http.Request) {}

// Show the page where the user can listen to an episode
func Listen(w http.ResponseWriter, r *http.Request) {}

// Show the page where a User who is a DM can upload an episode of a story
func Upload(w http.ResponseWriter, r *http.Request) {}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
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
