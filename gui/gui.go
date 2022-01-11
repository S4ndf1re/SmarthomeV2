package gui

import (
	"Smarthome/user"
	"Smarthome/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"sync"
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
	gui.cookieStore = sessions.NewCookieStore([]byte(os.Getenv(sessionEnvKey)))
	gui.setupLogin()
	_ = gui.addURLFunc(guiPath, gui.AuthorizeOrRedirect(gui.GuiHandle))
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
	container.AddToGui(apiMountPath, gui)
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
		gui.Containers[foundIdx].RemoveFromGui(apiMountPath, gui)
		gui.Containers = append(gui.Containers[:foundIdx], gui.Containers[foundIdx+1:]...)
	}
}

// AuthorizeOrRedirect tries to authorize the user session. If successful, the username and a nil error is returned.
// If the authorization is unsuccessful, the request is redirected to /login. An error is returned. After redirection (err != nil)
// The response should not get rewritten
func (gui *Gui) AuthorizeOrRedirect(callIfAuthorized func(string, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session, _ := gui.cookieStore.Get(r, sessionLogin)
		usernameInterface, ok := session.Values[sessionUsername]

		if !ok {
			http.Redirect(w, r, loginPath, http.StatusSeeOther)
			return
		}

		session.Options.MaxAge = sessionKeepAlive
		session.Options.Secure = true
		_ = session.Save(r, w)

		if username, ok := usernameInterface.(string); ok {
			callIfAuthorized(username, w, r)
		}
	}
}

// setupLogin registers all login utility as well as logout
func (gui *Gui) setupLogin() {
	err := gui.addURLFunc(loginApiPath, gui.LoginApi)
	util.LogIfErr("Gui.setupLogin()", err)

	err = gui.addURLFunc(loginPath, gui.Login)
	util.LogIfErr("Gui.setupLogin()", err)

	err = gui.addURLFunc(logoutPath, gui.Logout)
	util.LogIfErr("Gui.setupLogin()", err)
}

func authorize(username, password string) (*user.User, error) {
	refUser, err := user.Load(username)
	if err != nil {
		return nil, err
	}

	tryingUser := user.New(username, password)

	if refUser.Equals(tryingUser) {
		return refUser, nil
	}
	return nil, fmt.Errorf("password or username wrong")
}

// LoginApi logs in the user with the login form
func (gui *Gui) LoginApi(w http.ResponseWriter, r *http.Request) {
	session, _ := gui.cookieStore.Get(r, sessionLogin)
	if err := r.ParseForm(); err != nil {
		util.LogIfErr("Gui.Login()", err)
		return
	}

	username := r.PostForm.Get(sessionUsername)
	password := r.PostForm.Get(sessionPassword)

	registeredUser, err := authorize(username, password)
	if err != nil {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return
	}

	session.Values[sessionUsername] = registeredUser.Name
	session.Values[sessionPassword] = registeredUser.Password

	session.Options.MaxAge = sessionKeepAlive
	session.Options.Secure = false
	session.Options.SameSite = http.SameSiteLaxMode
	if err := session.Save(r, w); err != nil {
		util.LogIfErr("Gui.Login()", err)
		return
	}

	http.Redirect(w, r, emptyPath, http.StatusFound)
}

// Login serves the login page
func (gui *Gui) Login(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, loginPagePath)
}

// Logout removes the existing user session
func (gui *Gui) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := gui.cookieStore.Get(r, sessionLogin)
	session.Options.MaxAge = sessionDelete
	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, loginPath, http.StatusFound)
}

// GuiHandle handles the /gui requests. It serves the *Gui as json
func (gui *Gui) GuiHandle(_ string, w http.ResponseWriter, _ *http.Request) {
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
	http.FileServer(http.Dir(fileServeDirectory)).ServeHTTP(w, r)
}
