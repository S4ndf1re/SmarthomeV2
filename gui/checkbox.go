package gui

import (
	"fmt"
	"net/http"
	"sync"
)

// Checkbox is a simple On/Off State Checkbox with state change handlers
type Checkbox struct {
	Name              string `json:"name"`
	Text              string `json:"text"`
	OnOffStateRequest string `json:"onOffStateRequest"`
	OnOnStateRequest  string `json:"onOnStateRequest"`
	OnGetStateRequest string `json:"onGetStateRequest"`
	GuiType           string `json:"type"`

	getCurrent func(user string) bool
	onOnState  func(user string)
	onOffState func(user string)
	state      bool
	onChange   func(user string, state bool)
	mutex      sync.Mutex
}

// Creates a new default Checkbox with standard state tracking
func NewCheckbox(name, text string) *Checkbox {
	checkbox := new(Checkbox)
	checkbox.Text = text
	checkbox.Name = name
	checkbox.OnGetStateRequest = ""
	checkbox.OnOffStateRequest = ""
	checkbox.OnOnStateRequest = ""
	checkbox.GuiType = "gui.Checkbox"

	checkbox.getCurrent = func(user string) bool { return checkbox.state }
	checkbox.onOffState = func(user string) { checkbox.state = false }
	checkbox.onOnState = func(user string) { checkbox.state = true }
	checkbox.onChange = func(_ string, _ bool) {}
	checkbox.mutex = sync.Mutex{}
	return checkbox
}

// Type returns the Type "gui.Checkbox"
func (checkbox *Checkbox) Type() string {
	return checkbox.GuiType
}

// GetStatus returns the automatically managed state
// Note: If own state is handled, this function is rendered useless
func (checkbox *Checkbox) GetStatus() bool {
	return checkbox.state
}

// SetChangeCallback registers a callback that is called every time the internally managed state is changed
// callback is a function that receives the operating user and the new changed state
func (checkbox *Checkbox) SetChangeCallback(callback func(username string, state bool)) {
	checkbox.onChange = callback
}

// OverrideListeners overrides all internal state listeners. This overrides the internal state management.
// If the internal state management should be used again, use OverrideListeners in combination with DefaultListeners.
func (checkbox *Checkbox) OverrideListeners(onOffState func(string), onOnState func(string), onGetState func(string) bool) {
	checkbox.onOffState = onOffState
	checkbox.onOnState = onOnState
	checkbox.getCurrent = onGetState
}

// DefaultListeners returns the default internal state handlers
// The return order is OffState, OnState and GetState
func (checkbox *Checkbox) DefaultListeners() (func(string), func(string), func(string) bool) {
	return func(s string) {
			checkbox.state = false
		}, func(s string) {
			checkbox.state = true
		}, func(s string) bool {
			return checkbox.state
		}
}

// AddToGui adds all listeners an function callbacks to the *Gui
func (checkbox *Checkbox) AddToGui(mount string, gui *Gui) {
	checkbox.OnOnStateRequest = mount + checkbox.Name + "/onstate/click"
	_ = gui.addURLFunc(checkbox.OnOnStateRequest, func(w http.ResponseWriter, r *http.Request) {
		checkbox.mutex.Lock()
		defer checkbox.mutex.Unlock()
		username, err := gui.AuthorizeOrRedirect(w, r)
		if err != nil {
			return
		}
		checkbox.onOnState(username)
		checkbox.onChange(username, checkbox.getCurrent(username))
		_, _ = w.Write([]byte(fmt.Sprintf("{ \"status\": %t }", checkbox.getCurrent(username))))
	})

	checkbox.OnOffStateRequest = mount + checkbox.Name + "/offstate/click"
	_ = gui.addURLFunc(checkbox.OnOffStateRequest, func(w http.ResponseWriter, r *http.Request) {
		checkbox.mutex.Lock()
		defer checkbox.mutex.Unlock()
		username, err := gui.AuthorizeOrRedirect(w, r)
		if err != nil {
			return
		}
		checkbox.onOffState(username)
		checkbox.onChange(username, checkbox.getCurrent(username))
		_, _ = w.Write([]byte(fmt.Sprintf("{ \"status\": %t }", checkbox.getCurrent(username))))
	})

	checkbox.OnGetStateRequest = mount + checkbox.Name + "/state/get"
	_ = gui.addURLFunc(checkbox.OnGetStateRequest, func(w http.ResponseWriter, r *http.Request) {
		checkbox.mutex.Lock()
		defer checkbox.mutex.Unlock()
		username, err := gui.AuthorizeOrRedirect(w, r)
		if err != nil {
			return
		}
		checkbox.getCurrent(username)
		_, _ = w.Write([]byte(fmt.Sprintf("{ \"status\": %t }", checkbox.getCurrent(username))))
	})
}

// RemoveFromGui removes all handlers from the *Gui
func (checkbox *Checkbox) RemoveFromGui(mount string, gui *Gui) {
	checkbox.OnOnStateRequest = mount + checkbox.Name + "/onstate/click"
	checkbox.OnOffStateRequest = mount + checkbox.Name + "/offstate/click"
	checkbox.OnGetStateRequest = mount + checkbox.Name + "/state/get"

	_ = gui.removeURLFunc(checkbox.OnOnStateRequest)
	_ = gui.removeURLFunc(checkbox.OnOffStateRequest)
	_ = gui.removeURLFunc(checkbox.OnGetStateRequest)
}
