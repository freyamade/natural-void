package naturalvoid

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/go-ldap/ldap"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type message struct {
	Type string
	Text string
}

var styleTheme = map[string]string{
	"dread": "#2C0047",
}

func NewRouter() chi.Router {
	conf := GetConf()
	r := chi.NewRouter()

	// Add middleware
	r.Use(ensureTrailingSlash)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Logger)
	r.Use(conf.CSRF)
	r.Use(context.ClearHandler)

	// Register the routes
	r.Get("/", Index)
	r.Get("/manifest/", Manifest)
	r.Get("/login/", LoginForm)
	r.Post("/login/", Login)
	r.Get("/logout/", Logout)
	r.Get("/story/{story:\\d+}/", Episodes)
	r.Get("/episode/{story:\\d+}/{episode:\\d+}/", Listen)
	r.Post("/episode/{story:\\d+}/{episode:\\d+}/delete/", DeleteEpisode)
	r.Get("/upload/", UploadForm)
	r.Post("/upload/", UploadEpisode)
	// Serve the static files
	fileServer(r, "/static/", http.Dir("./static"))
	// Also serve the episodes
	fileServer(r, "/episodes/", http.Dir("./episodes"))
	return r
}

// Define Route Handlers here
func Index(w http.ResponseWriter, r *http.Request) {
	// Get the list of Stories from the DB
	data := map[string]interface{}{}
	dao, err := GetDAO()
	if err != nil {
		fmt.Println(err)
		conf := GetConf()
		session, _ := conf.SessionStore.Get(r, "session")
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
	} else {
		stories := []Story{}
		dao.DB.Find(&stories)
		data["Stories"] = stories
	}
	render(w, r, "index.tmpl", data)
}

// Display a form to the User to allow them to log in
func LoginForm(w http.ResponseWriter, r *http.Request) {
	// Check to make sure the user hasn't already logged in
	conf := GetConf()
	session, _ := conf.SessionStore.Get(r, "session")
	if session.Values["Authenticated"] == true {
		// Redirect back to the index with a message saying they logged in.
		session.AddFlash("success:You are already logged in!")
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)
		return
	}
	// Get the CSRF token
	data := map[string]interface{}{
		"CSRF":  csrf.TemplateField(r),
		"Title": "Login",
	}
	render(w, r, "login.tmpl", data)
}

// Handle logging in of a user by checking against LDAP
func Login(w http.ResponseWriter, r *http.Request) {
	// Store important things in the session
	data := map[string]interface{}{
		"CSRF":     csrf.TemplateField(r),
		"Title":    "Login",
	}
	conf := GetConf()
	session, _ := conf.SessionStore.Get(r, "session")
	// Attempt to auth the user
	// Attempt to parse the form
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		session.AddFlash("warning:Could not parse form. Try again later.")
		session.Save(r, w)
		render(w, r, "login.tmpl", data)
		return
	}
	// Attempt to login to the LDAP server
	sentUsername := r.PostForm.Get("username")
	sentPassword := r.PostForm.Get("password")
	data["Username"] = sentUsername

	// Validate the sent username and password
	ok, err := isValidAuth(sentUsername, sentPassword)
	if err != nil {
		fmt.Println(err)
		session.AddFlash("warning:Error during LDAP verification. Please try again later.")
		session.Save(r, w)
		render(w, r, "login.tmpl", data)
		return
	}

	if ok {
		session.Values["Authenticated"] = true
		session.Values["Username"] = sentUsername
		session.AddFlash("success:You have logged in successfully!")
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)
	} else {
		// Re-render the login form with an error message
		session.AddFlash("danger:Invalid username or password. Please check your details and try again.")
		session.Save(r, w)
		data["Username"] = sentUsername
		render(w, r, "login.tmpl", data)
	}
}

// Handle logging out of a logged in user
func Logout(w http.ResponseWriter, r *http.Request) {
	// Remove the authenticated flag from the session
	conf := GetConf()
	session, _ := conf.SessionStore.Get(r, "session")
	if session.Values["Authenticated"] == true {
		session.Values["Authenticated"] = false
		session.Values["Username"] = ""
		session.AddFlash("success:You have been logged out successfully!")
	} else {
		// Redirect back to the index with a message saying they logged in.
		session.AddFlash("danger:You have to be logged in to log out!")
	}
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
}

// Show the list of episodes in order of newest first
func Episodes(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"CSRF": csrf.TemplateField(r),
	}

	dao, err := GetDAO()
	if err != nil {
		fmt.Println(err)
		conf := GetConf()
		session, _ := conf.SessionStore.Get(r, "session")
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		data["Title"] = "Episodes"
	} else {
		storyID := chi.URLParam(r, "story")
		st := Story{}
		user := User{}
		episodes := []Episode{}

		dao.DB.Find(&st, storyID).Related(&user)
		dao.DB.Order("number DESC").Find(&episodes)
		conf := GetConf()
		session, _ := conf.SessionStore.Get(r, "session")

		data["Title"] = fmt.Sprintf("%s Episodes", st.Name)
		data["Story"] = st
		data["Episodes"] = episodes
		data["IsOwner"] = session.Values["Username"] == user.Username
	}
	render(w, r, "episode_list.tmpl", data)
}

// Show the page where the user can listen to an episode
func Listen(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}

	dao, err := GetDAO()
	if err != nil {
		fmt.Println(err)
		conf := GetConf()
		session, _ := conf.SessionStore.Get(r, "session")
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		data["Title"] = "Listen to Natural Void"
	} else {
		st := Story{}
		ep := Episode{}
		storyID := chi.URLParam(r, "story")
		episode := chi.URLParam(r, "episode")
		dao.DB.Find(&st, storyID)
		dao.DB.Where("story_id = ? AND number = ?", storyID, episode).Find(&ep)

		// Get next and previous episodes if they exist
		prev := Episode{}
		next := Episode{}
		dao.DB.Where("number = ? AND story_id = ?", (ep.Number - 1), storyID).First(&prev)
		dao.DB.Where("number = ? AND story_id = ?", (ep.Number + 1), storyID).First(&next)

		data["Title"] = fmt.Sprintf("Listen to %s", ep.Name)
		data["Episode"] = ep
		data["Story"] = st
		data["Prev"] = prev
		data["Next"] = next
	}
	render(w, r, "listen.tmpl", data)
}

// Show the page where a User who is a DM can upload an episode of a story
func UploadForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"CSRF":  csrf.TemplateField(r),
		"Title": "Upload an episode",
	}

	conf := GetConf()
	session, _ := conf.SessionStore.Get(r, "session")
	isDM, err := isDM(session.Values["Username"].(string))
	if err != nil {
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	if !(session.Values["Authenticated"] == true) && !(isDM) {
		// Redirect to the index with an error message
		session.AddFlash("danger:You must be logged in and be running a story to access this page.")
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)
		return
	}

	stories, err := getStories(session.Values["Username"].(string))
	if err != nil {
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	// Display a form allowing the user to upload an episode
	data["Stories"] = stories
	render(w, r, "upload.tmpl", data)
}

// Handle the uploading of an Episode into the DB
func UploadEpisode(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"CSRF":  csrf.TemplateField(r),
		"Title": "Upload an episode",
	}

	conf := GetConf()
	session, _ := conf.SessionStore.Get(r, "session")
	isDM, err := isDM(session.Values["Username"].(string))
	if err != nil {
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	if !(session.Values["Authenticated"] == true) && !(isDM) {
		// Redirect to the index with an error message
		session.AddFlash("danger:You must be logged in and be running a story to access this page.")
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)
		return
	}

	// Get the user's stories
	stories, err := getStories(session.Values["Username"].(string))
	if err != nil {
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}
	data["Stories"] = stories

	// Parse the form and check for valid params
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err)
		session.AddFlash("warning:Error occurred when parsing the form, try again later!")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}
	data["Name"] = r.PostForm.Get("name")
	data["Description"] = r.PostForm.Get("description")

	// Attempt to connect to the DB
	dao, err := GetDAO()
	if err != nil {
		fmt.Println(err)
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	// Ensure all fields have been sent
	if r.PostForm.Get("name") == "" || r.PostForm.Get("description") == "" || r.PostForm.Get("story") == "" {
		session.AddFlash("danger:All of the text fields must be filled in.")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	// Validate the story id
	st := Story{}
	err = dao.DB.Find(&st, r.PostForm.Get("story")).Error
	if err != nil {
		session.AddFlash("danger:Invalid Story Id.")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	// Check that no other episodes have the same name in the story
	var count int
	dao.DB.Model(&Episode{}).Where("story_id = ? AND name = ?", st.ID, r.PostForm.Get("name")).Count(&count)
	if count > 0 {
		session.AddFlash("danger:An episode with the specified name already exists for the chosen story.")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	// Check that the file has been sent
	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		session.AddFlash("danger:Error reading file. Make sure it has been sent correctly.")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	// Check mime type of file to ensure it is audio
	contentType := handler.Header["Content-Type"][0]
	if strings.Split(contentType, "/")[0] != "audio" {
		session.AddFlash("danger:Please ensure the uploaded file is an audio file.")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}

	// If we get to this point, it's safe to say that the upload is good

	// Calculate the number for this episode
	var number int
	dao.DB.Model(&Episode{}).Where("story_id = ?", st.ID).Count(&number)
	number += 1
	// Create the new model
	ep := Episode{
		Name:        r.PostForm.Get("name"),
		Description: strings.Split(r.PostForm.Get("description"), "\n"),
		Number:      number,
		StoryID:     st.ID,
	}

	// Now put the file in the correct directory `./episodes/{storyID}/{episodeID}`
	// Ensure the dir exists
	os.MkdirAll(fmt.Sprintf("./episodes/%d/", st.ID), 0755)
	storeFile, err := os.OpenFile(fmt.Sprintf("./episodes/%d/%d", st.ID, number), os.O_WRONLY|os.O_CREATE, 0644)
	defer storeFile.Close()
	if err != nil {
		session.AddFlash("danger:An error occurred when trying to upload the file. Please try again later.")
		session.Save(r, w)
		render(w, r, "upload.tmpl", data)
		return
	}
	io.Copy(storeFile, file)
	dao.DB.Create(&ep)

	session.AddFlash(fmt.Sprintf("success:Successfully uploaded episode %d of %s.", number, st.Name))
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/story/%d/", st.ID), 303)
}

// Delete an episode if the current auth'd user owns it
func DeleteEpisode(w http.ResponseWriter, r *http.Request) {
	conf := GetConf()
	session, _ := conf.SessionStore.Get(r, "session")
	storyID := chi.URLParam(r, "story")
	// Attempt to connect to the DB
	dao, err := GetDAO()
	if err != nil {
		fmt.Println(err)
		session.AddFlash("warning:Could not connect to database!")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/story/%d/", storyID), 303)
		return
	}

	user := User{}
	st := Story{}
	ep := Episode{}
	episode := chi.URLParam(r, "episode")

	// Delete the episode from the DB and also delete the episode file from the file system
	err = dao.DB.Find(&st, storyID).Related(&user).Error
	if err != nil {
		session.AddFlash("danger:Invalid Story ID.")
		session.Save(r, w)
		http.Redirect(w, r, "/", 303)
		return
	}

	// Ensure the episode exists
	err = dao.DB.Where("story_id = ? AND number = ?", storyID, episode).Find(&ep).Error
	if err != nil {
		session.AddFlash("danger:Invalid Episode Number.")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/story/%d/", st.ID), 303)
		return
	}

	// Ensure the requester owns the story
	if session.Values["Authenticated"] != true || session.Values["Username"] != user.Username {
		session.AddFlash("danger:Episodes can only be deleted by the owner of the story.")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/story/%d/", st.ID), 303)
		return
	}

	// Delete from the file system
	path := fmt.Sprintf("./episodes/%d/%d", st.ID, ep.Number)
	err = os.Remove(path)
	if err != nil {
		fmt.Println(err)
		session.AddFlash("danger:Couldn't delete episode file from filesystem.")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/story/%d/", st.ID), 303)
		return
	}
	// Delete from DB
	dao.DB.Delete(&ep)
	session.AddFlash(fmt.Sprintf("success:Successfully deleted episode %d of %s.", ep.Number, st.Name))
	session.Save(r, w)
	http.Redirect(w, r, fmt.Sprintf("/story/%d/", st.ID), 303)
}

// Generate the manifest
func Manifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	render(w, r, "manifest.tmpl", map[string]interface{}{})
}

// HELPERS

// Helper to ensure necessary data is always passed to Template
func render(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	// Add session data to the map
	conf := GetConf()
	session, _ := conf.SessionStore.Get(r, "session")

	// Add necessary data to data, and then render the specified template
	var style string
	if session.Values["style"] == nil {
		style = "dread"
	} else {
		style = session.Values["style"].(string)
	}
	session.Values["style"] = style
	data["Style"] = style
	data["Theme"] = styleTheme[style]
	data["Session"] = session.Values

	// Check whether or not the current user is a DM for a story
	if session.Values["Authenticated"] == true {
		// Check if they have any stories
		isDM, _ := isDM(session.Values["Username"].(string)) // Ignore error here in case we double post the same message
		data["IsDM"] = isDM
	}

	// Check flash messages
	if flashes := session.Flashes(); len(flashes) > 0 {
		var messages []message
		if data["Messages"] != nil {
			messages = data["Messages"].([]message)
		} else {
			messages = []message{}
		}
		for _, msg := range flashes {
			splitMsg := strings.Split(msg.(string), ":")
			// Split the type and message from the message
			messages = append(messages, message{Type: splitMsg[0], Text: splitMsg[1]})
		}
		data["Messages"] = messages
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

// isDM checks the Database for stories owned by the user with the given username
func isDM(username string) (bool, error) {
	user := User{}
	dao, err := GetDAO()
	if err != nil {
		return false, err
	}
	dao.DB.Where("username = ?", username).Find(&user)
	var count uint
	dao.DB.Model(&Story{}).Where("user_id = ?", user.ID).Count(&count)
	return count > 0, nil
}

// Returns a slice of Story structs owned by the User with the supplied username
func getStories(username string) ([]Story, error) {
	user := User{}
	dao, err := GetDAO()
	if err != nil {
		return nil, err
	}
	dao.DB.Where("username = ?", username).Find(&user)
	stories := []Story{}
	dao.DB.Where("user_id = ?", user.ID).Find(&stories)
	return stories, nil
}

// Check the passed username and password against the LDAP DB
func isValidAuth(username string, password string) (bool, error) {
	conf := GetConf()
	// Connect to the LDAP server
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", conf.LDAPHost, conf.LDAPPort))
	if err != nil {
		return false, err
	}
	defer l.Close()

	// Bind as the user in the conf
	err = l.Bind(conf.LDAPUser, conf.LDAPPass)
	if err != nil {
		return false, err
	}

	// Search for the given username
	filter := fmt.Sprintf("(&(cn=%s))", username)
	searchRequest := ldap.NewSearchRequest("dc=naturalvoid,dc=com", ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, []string{"dn"}, nil)

	search, err := l.Search(searchRequest)
	if err != nil {
		return false, err
	}

	// Check that only one entry has been returned
	if len(search.Entries) != 1 {
		return false, errors.New("User either does not exist or returned more than 1 LDAP entry")
	}

	// Get the DN of the user and attempt to bind as that user to see if the password is valid
	userDN := search.Entries[0].DN
	err = l.Bind(userDN, password)
	if err != nil {
		return false, err
	}

	// If we make it here, the auth is valid
	return true, nil
}

