package naturalvoid

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	r.Use(middleware.DefaultCompress)
	// Register the routes
	r.Get("/", Index)
	r.Get("/manifest/", Manifest)
	r.Get("/login/", LoginForm)
	r.Post("/login/", Login)

	r.Get("/listen/{episode:\\d+}/", Listen)
	// Serve the static files
	fileServer(r, "/static/", http.Dir("./static"))
	// Also serve the episodes
	fileServer(r, "/episodes/", http.Dir("./episodes"))
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
	render(w, "index.tmpl", data)
}

// Display a form to the User to allow them to log in
func LoginForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Login",
	}
	render(w, "login.tmpl", data)
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
	render(w, "login.tmpl", data)
}

// Handle logging out of a logged in user
func Logout(w http.ResponseWriter, r *http.Request) {}

// Show the list of episodes in order of newest first
func Episodes(w http.ResponseWriter, r *http.Request) {}

// Show the page where the user can listen to an episode
func Listen(w http.ResponseWriter, r *http.Request) {
	st := Story{}
	ep := Episode{}
	episodeID := chi.URLParam(r, "episode")
	dao := GetDAO()
	dao.DB.Find(&ep, episodeID).Related(&st)

	// Get next and previous episodes if they exist
	prev := Episode{}
	next := Episode{}
	dao.DB.Where("number = ? AND story_id = ?", (ep.Number - 1), st.ID).First(&prev)
	dao.DB.Where("number = ? AND story_id = ?", (ep.Number + 1), st.ID).First(&next)

	data := map[string]interface{}{
		"Title":   "Listen",
		"Episode": ep,
		"Story":   st,
		"Prev": prev,
		"Next": next,
	}
	render(w, "listen.tmpl", data)
}

// Show the page where a User who is a DM can upload an episode of a story
func UploadForm(w http.ResponseWriter, r *http.Request) {}

// Handle the uploading of an Episode into the DB
func UploadEpisode(w http.ResponseWriter, r *http.Request) {}

// Generate the manifest
func Manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	render(w, "manifest.tmpl", map[string]interface{}{})
}

// HELPERS

// Helper to ensure necessary data is always passed to Template
func render(w http.ResponseWriter, name string, data map[string]interface{}) {
	// Add necessary data to data, and then render the specified template
	style := "dread"  // Overwrite with logic to get style from session
	// Map the style name to its theme colour
	styleTheme := map[string]string {
		"dread": "#2C0047",
	}
	data["Style"] = style
	data["Theme"] = styleTheme[style]

	// Generate and parse the templates
	tmpl, err := template.ParseFiles("templates/layout.tmpl", fmt.Sprintf("templates/%s", name))
	if err != nil {
		panic(err)
	}
	err = tmpl.ExecuteTemplate(w, name, &data)
	if err != nil {
		panic(err)
	}
}

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
