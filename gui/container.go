package gui

import (
	"Smarthome/util"
	"net/http"
	"sync"
)

// Container is a container that contains a list of children. When clicked in the gui, onInit is called
type Container struct {
	Name            string `json:"name"`
	Text            string `json:"text"`
	OnInitRequest   string `json:"onInitRequest"`
	onInit          func(user string)
	OnUnloadRequest string `json:"onUnloadRequest"`
	onUnload        func(user string)
	List            []Child `json:"list"`
	mutex           sync.Mutex
}

// NewContainer creates a new container and initializes it
func NewContainer(name string, text string, onInit func(string), onUnload func(string)) *Container {
	container := new(Container)
	container.Name = name
	container.Text = text
	container.OnInitRequest = ""
	container.onInit = onInit
	container.OnUnloadRequest = ""
	container.onUnload = onUnload
	container.List = make([]Child, 0)
	container.mutex = sync.Mutex{}
	return container
}

// Add adds a child to the container
// The new child is appended at the end
func (container *Container) Add(child Child) {
	container.mutex.Lock()
	defer container.mutex.Unlock()
	container.List = append(container.List, child)
}

func (container *Container) handleOnInitRequest(username string, _ http.ResponseWriter, _ *http.Request) {
	container.onInit(username)
}

func (container *Container) handleOnUnloadRequest(username string, _ http.ResponseWriter, _ *http.Request) {
	container.onUnload(username)
}

// AddToGui adds the container and all its children to the *Gui
func (container *Container) AddToGui(mount string, gui *Gui) {
	container.mutex.Lock()
	defer container.mutex.Unlock()
	newMount := mount + container.Name + pathSeparator

	for _, child := range container.List {
		child.AddToGui(newMount, gui)
	}

	container.OnInitRequest = mount + container.Name + containerInitPath
	err := gui.addURLFunc(container.OnInitRequest, gui.AuthorizeOrRedirect(container.handleOnInitRequest))
	util.LogIfErr("Container.AddToGui()", err)

	container.OnUnloadRequest = mount + container.Name + containerUnloadPath
	err = gui.addURLFunc(container.OnUnloadRequest, gui.AuthorizeOrRedirect(container.handleOnUnloadRequest))
	util.LogIfErr("Container.AddToGui()", err)
}

// RemoveFromGui removes the container and all its children from the *Gui
func (container *Container) RemoveFromGui(mount string, gui *Gui) {
	container.mutex.Lock()
	defer container.mutex.Unlock()
	newMount := mount + container.Name + pathSeparator

	for _, child := range container.List {
		child.RemoveFromGui(newMount, gui)
	}

	container.OnInitRequest = mount + container.Name + containerInitPath
	err := gui.removeURLFunc(container.OnInitRequest)
	util.LogIfErr("Container.RemoveFromGui()", err)

	container.OnUnloadRequest = mount + container.Name + containerUnloadPath
	err = gui.removeURLFunc(container.OnUnloadRequest)
	util.LogIfErr("Container.RemoveFromGui()", err)
}
