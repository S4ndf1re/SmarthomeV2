package gui

import (
	"Smarthome/util"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// Data is a websocket connection to the gui with one Child that is the updatable data
type Data struct {
	Name            string `json:"name"`
	data            Child
	dataMount       string
	dataGui         *Gui
	mutex           sync.Mutex
	UpdateRequest   string `json:"updateRequest"`
	UpdateSocket    string `json:"updateSocket"`
	updateFunctions map[string]func(Child)
	GuiType         string `json:"type"`
}

// NewData creates a new Data struct
func NewData(name string, initialData Child) *Data {
	data := new(Data)
	data.Name = name
	data.data = initialData
	data.mutex = sync.Mutex{}
	data.updateFunctions = make(map[string]func(Child))
	data.GuiType = DataType
	return data
}

// Type returns the type "gui.Data"
func (data *Data) Type() string {
	return data.GuiType
}

// Update can be called to update the underlying Child data.
// After update, all updateFunctions are called. Hence all Websocket connections are updated
func (data *Data) Update(newData Child) {
	data.mutex.Lock()
	defer data.mutex.Unlock()
	// Remove old paths. This line makes buttons etc clean up correctly
	data.data.RemoveFromGui(data.dataMount, data.dataGui)

	data.data = newData

	// Add new paths to gui where data is registered. This makes buttons etc. work
	data.data.AddToGui(data.dataMount, data.dataGui)
	for _, f := range data.updateFunctions {
		f(data.data)
	}
}

// addUpdateFunction adds a callback for when the update is triggered.
// A unique id is returned. The id can be used to remove the callback from the list by calling removeUpdateFunction
func (data *Data) addUpdateFunction(updateFunction func(Child)) string {
	data.mutex.Lock()
	defer data.mutex.Unlock()
	var ident string
	// Make sure the key is unique. Because the key is large, it should only take one or two tries max
	for {
		ident = util.RandomBase64Bytes(128)
		if _, ok := data.updateFunctions[ident]; !ok {
			break
		}
	}
	data.updateFunctions[ident] = updateFunction

	return ident
}

// removeUpdateFunction removes a update function from the internal map.
// The [ident] parameter is the value returned by addUpdateFunction
func (data *Data) removeUpdateFunction(ident string) {
	data.mutex.Lock()
	defer data.mutex.Unlock()
	delete(data.updateFunctions, ident)
}

// handleRequest handles the simple get request on the *Data
func (data *Data) handleRequest(_ string, w http.ResponseWriter, _ *http.Request) {
	data.mutex.Lock()
	defer data.mutex.Unlock()

	buffer, err := json.Marshal(data.data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(buffer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleSocket handles all websocket requests to *Data
func (data *Data) handleSocket(_ string, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateCallback := func(child Child) {
		buffer, err := json.Marshal(data.data)
		if err != nil {
			return
		}

		if err = conn.WriteMessage(websocket.TextMessage, buffer); err != nil {
			return
		}
	}

	callbackHandle := data.addUpdateFunction(updateCallback)
	defer data.removeUpdateFunction(callbackHandle)

	// Block until close
	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if messageType == websocket.CloseMessage {
			break
		}
	}
}

// AddToGui registers both websocket and get request handlers to the *Gui
func (data *Data) AddToGui(mount string, gui *Gui) {
	data.UpdateRequest = mount + data.Name + dataRequestPath
	data.UpdateSocket = mount + data.Name + dataSocketPath

	err := gui.addURLFunc(data.UpdateRequest, gui.AuthorizeOrRedirect(data.handleRequest))
	util.LogIfErr("Data.AddToGui()", err)

	err = gui.addURLFunc(data.UpdateSocket, gui.AuthorizeOrRedirect(data.handleSocket))
	util.LogIfErr("Data.AddToGui()", err)

	data.dataMount = mount
	data.dataGui = gui
	data.data.AddToGui(mount, gui)
}

// RemoveFromGui removes all registered handlers from the *Gui
func (data *Data) RemoveFromGui(mount string, gui *Gui) {
	data.UpdateRequest = mount + data.Name + dataRequestPath
	data.UpdateSocket = mount + data.Name + dataSocketPath

	err := gui.removeURLFunc(data.UpdateRequest)
	util.LogIfErr("Data.RemoveFromGui()", err)

	err = gui.removeURLFunc(data.UpdateSocket)
	util.LogIfErr("Data.RemoveFromGui()", err)

	data.data.RemoveFromGui(data.dataMount, data.dataGui)
}
