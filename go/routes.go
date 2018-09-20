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
	r.Use(middleware.Logger)
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
	render(w, r, "index.tmpl", data)
}

// Display a form to the User to allow them to log in
func LoginForm(w http.ResponseWriter, r *http.Request) {
    // Check to make sure the user hasn't already logged in
    conf := GetConf()
    session, _ := conf.SessionStore.Get(r, "session")
    fmt.Println(session.Values)
    if session.Values["Authenticated"] == true {
        // Redirect back to the index with a message saying they logged in.
        session.AddFlash("You are already logged in!")
        session.Save(r, w)
        http.Redirect(w, r, "/", 303)
        return
    }
	data := map[string]interface{}{
		"Title": "Login",
	}
	render(w, r, "login.tmpl", data)
}

// Handle logging in of a user by checking against LDAP
func Login(w http.ResponseWriter, r *http.Request) {
    // Attempt to auth the user
	// Attempt to parse the form
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
    // Create some sample login data for now until I can add LDAP later
    auth := map[string]string {
        "username": "crnbrdrck",
        "password": "password",
    }
    sentUsername := r.Form.Get("username")
    sentPassword := r.Form.Get("password")

    // Validate the sent username and password
    if sentUsername == auth["username"] && sentPassword == auth["password"] {
        // Store important things in the session
        conf := GetConf()
        session, _ := conf.SessionStore.Get(r, "session")
        session.Values["Authenticated"] = true
        session.Values["Username"] = sentUsername
        session.AddFlash("You have logged in successfully!")
        session.Save(r, w)
        http.Redirect(w, r, "/", 303)
    }
}

// Handle logging out of a logged in user
func Logout(w http.ResponseWriter, r *http.Request) {
    // Remove the authenticated flag from the session
}

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
	render(w, r, "listen.tmpl", data)
}

// Show the page where a User who is a DM can upload an episode of a story
func UploadForm(w http.ResponseWriter, r *http.Request) {}

// Handle the uploading of an Episode into the DB
func UploadEpisode(w http.ResponseWriter, r *http.Request) {}

// Generate the manifest
func Manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	render(w, r, "manifest.tmpl", map[string]interface{}{})
}

// HELPERS

// Helper to ensure necessary data is always passed to Template
func render(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	// Add necessary data to data, and then render the specified template
	style := "dread"  // Overwrite with logic to get style from session
	// Map the style name to its theme colour
	styleTheme := map[string]string {
		"dread": "#2C0047",
	}
	data["Style"] = style
	data["Theme"] = styleTheme[style]

    // Add session data to the map
    conf := GetConf()
    session, _ := conf.SessionStore.Get(r, "session")
    data["Session"] = session.Values

    // Check flash messages
    fmt.Println("Checking flashes")
    if flashes := session.Flashes(); len(flashes) > 0 {
        fmt.Println("Generating messages")
        var messages []message
        if data["Messages"] != nil {
            messages = data["Messages"].([]message)
        } else {
            messages = []message{}
        }
        for _, msg := range flashes {
            messages = append(messages, message{Type: "success", Text: msg.(string)})
        }
        data["Messages"] = messages
        fmt.Println(messages)
    }
    session.Save(r, w)

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
		r.Get(path, http.RedirectHandler(path+"/", 303).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
