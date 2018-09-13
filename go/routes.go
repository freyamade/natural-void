package naturalvoid

import (
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"strings"
	// "github.com/gorilla/sessions"
)

type message struct {
	Type string
	Text string
}

func NewRouter() chi.Router {
	r := chi.NewRouter()
	// Add middleware
	r.Use(ensureTrailingSlash)
	// Register the routes
	r.Get("/", Index)
	r.Get("/login/", LoginForm)
	r.Post("/login/", Login)
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
	data := map[string]interface{}{
		"Stories": stories,
	}
	// Generate and parse the templates
	tmpl, err := template.ParseFiles("templates/layout.tmpl", "templates/index.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(w, "index.tmpl", &data)
	if err != nil {
		panic(err)
	}
}

// Display a form to the User to allow them to log in
func LoginForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login",
	}
	// Generate and parse the templates
	tmpl, err := template.ParseFiles("templates/layout.tmpl", "templates/login.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(w, "login.tmpl", data)
	if err != nil {
		panic(err)
	}
}

// Handle logging in of a user by checking against LDAP
func Login(w http.ResponseWriter, r *http.Request) {
	// For now, just render the login form with the passed username
	// Attempt to parse the form
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{
		"Title":    "Login",
		"Username": r.Form.Get("username"),
	}
	// Generate and parse the templates
	tmpl, err := template.ParseFiles("templates/layout.tmpl", "templates/login.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(w, "login.tmpl", data)
	if err != nil {
		panic(err)
	}
}

// Handle logging out of a logged in user
func Logout(w http.ResponseWriter, r *http.Request) {}

// Show the list of episodes in order of newest first
func Episodes(w http.ResponseWriter, r *http.Request) {}

// Show the page where the user can listen to an episode
func Listen(w http.ResponseWriter, r *http.Request) {}

// Show the page where a User who is a DM can upload an episode of a story
func UploadForm(w http.ResponseWriter, r *http.Request) {}

// Handle the uploading of an Episode into the DB
func UploadEpisode(w http.ResponseWriter, r *http.Request) {}

// HELPERS

// Middleware to add a trailing slash to the url if one is missing
func ensureTrailingSlash(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var path string
		ctx := chi.RouteContext(r.Context())
		if ctx.RoutePath != "" {
			path = ctx.RoutePath
		} else {
			path = r.URL.Path
		}
		if len(path) > 1 && path[len(path)-1] != '/' {
			ctx.RoutePath = path + "/"
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

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
