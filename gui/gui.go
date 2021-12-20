package gui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	sessionLogin    = "session_login"
	sessionUsername = "username"
	sessionPassword = "password"
)

// Gui represents a web ui
type Gui struct {
	Containers  []*Container `json:"containers"`
	urlList     map[string]func(w http.ResponseWriter, r *http.Request)
	cookieStore *sessions.CookieStore
	mutex       sync.Mutex
}

// NewGui creates a new, empty *Gui
func NewGui() *Gui {
	gui := new(Gui)
	gui.Containers = make([]*Container, 0)
	gui.urlList = make(map[string]func(w http.ResponseWriter, r *http.Request))
	gui.cookieStore = sessions.NewCookieStore([]byte(os.Getenv("COOKIE_KEY")))
	gui.setupLogin()
	_ = gui.addURLFunc("/gui", gui.GuiHandle)
	gui.mutex = sync.Mutex{}
	return gui
}

// addURLFunc adds a new http callback on given [path] to the gui.
// Returns an error if the path is already registered
func (gui *Gui) addURLFunc(path string, callback func(w http.ResponseWriter, r *http.Request)) error {
	_, err := gui.getURLFunc(path)
	if err != nil {
		gui.urlList[path] = callback
		return nil
	}
	return fmt.Errorf("%s is already registered", path)
}

// removeURLFunc removes the http callback for given [path].
// Returns an error if the path is not registered
func (gui *Gui) removeURLFunc(path string) error {
	if _, err := gui.getURLFunc(path); err != nil {
		return err
	}
	delete(gui.urlList, path)
	return nil
}

// getURLFunc finds and returns a http callback for given [path].
// Returns nil and an error if the path is not registered.
// If the path is registered, the http callback is returned
func (gui *Gui) getURLFunc(path string) (func(w http.ResponseWriter, r *http.Request), error) {
	gui.mutex.Lock()
	defer gui.mutex.Unlock()
	val, ok := gui.urlList[path]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s not found", path))
	}
	return val, nil
}

// ServeHTTP serves all go http requests
func (gui *Gui) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	callbackFunction, err := gui.getURLFunc(r.URL.Path)
	if err != nil {
		gui.HandleRoot(w, r)
		return
	}
	callbackFunction(w, r)
}

// AddContainer adds the [container] to the container list of the gui
// It also registers all components to the handler
func (gui *Gui) AddContainer(container *Container) {
	gui.Containers = append(gui.Containers, container)
	container.AddToGui("/api/", gui)
}

// RemoveContainer removes the [container] from the gui container list.
// The removal is done by name.
// All registered gui components will get removed as well.
func (gui *Gui) RemoveContainer(container *Container) {
	foundIdx := -1
	for idx, value := range gui.Containers {
		if value.Name == container.Name {
			foundIdx = idx
		}
	}

	if foundIdx >= 0 {
		gui.Containers[foundIdx].RemoveFromGui("/api/", gui)
		gui.Containers = append(gui.Containers[:foundIdx], gui.Containers[foundIdx+1:]...)
	}
}

// AuthorizeOrRedirect tries to authorize the user session. If successful, the username and a nil error is returned.
// If the authorization is unsuccessful, the request is redirected to /login. An error is returned. After redirection (err != nil)
// The response should not get rewritten
func (gui *Gui) AuthorizeOrRedirect(w http.ResponseWriter, r *http.Request) (string, error) {
	session, _ := gui.cookieStore.Get(r, sessionLogin)
	usernameInterface, ok := session.Values[sessionUsername]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return "", fmt.Errorf("no username found. Redirecting to /login")
	}
	if username, ok := usernameInterface.(string); ok {
		return username, nil
	} else {
		return "", fmt.Errorf("cookie does not contain string")
	}
}

// setupLogin registers all login utility as well as logout
func (gui *Gui) setupLogin() {
	if err := gui.addURLFunc("/util/login", gui.LoginApi); err != nil {
		log.Printf("setupLogin(): %s\n", err)
	}
	if err := gui.addURLFunc("/login", gui.Login); err != nil {
		log.Printf("setupLogin(): %s\n", err)
	}
	if err := gui.addURLFunc("/logout", gui.Logout); err != nil {
		log.Printf("setupLogin(): %s\n", err)
	}
}

// LoginApi logs in the user with the login form
func (gui *Gui) LoginApi(w http.ResponseWriter, r *http.Request) {
	session, _ := gui.cookieStore.Get(r, sessionLogin)
	if err := r.ParseForm(); err != nil {
		log.Printf("Login(): %s\n", err)
		return
	}

	session.Values[sessionUsername] = r.PostForm.Get(sessionUsername)
	session.Values[sessionPassword] = r.PostForm.Get(sessionPassword)

	if err := session.Save(r, w); err != nil {
		log.Printf("Login(): %s\n", err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Login serves the login page
func (gui *Gui) Login(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/login.html")
}

// Logout removes the existing user session
func (gui *Gui) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := gui.cookieStore.Get(r, sessionLogin)
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// GuiHandle handles the /gui requests. It serves the *Gui as json
func (gui *Gui) GuiHandle(w http.ResponseWriter, r *http.Request) {

	if _, err := gui.AuthorizeOrRedirect(w, r); err != nil {
		return
	}

	data, err := json.Marshal(gui)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HandleRoot handles all other requests that may be part of the /html folder
func (gui *Gui) HandleRoot(w http.ResponseWriter, r *http.Request) {
	if _, err := gui.AuthorizeOrRedirect(w, r); err != nil {
		return
	}
	http.FileServer(http.Dir("html/")).ServeHTTP(w, r)
}
