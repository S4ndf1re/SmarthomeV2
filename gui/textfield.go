package gui

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type textInputData struct {
	Text string `json:"text"`
}

// Textfield is a textinput component
type Textfield struct {
	Name           string `json:"name"`
	Text           string `json:"text"`
	UpdateRequest  string `json:"updateRequest"`
	GuiType        string `json:"type"`
	onChange       func(username string, text string)
	currentContent string
}

// NewTextfield constructs a new text field with a default change listener
// Internal state is kept track of
func NewTextfield(name string, text string, onChange func(username string, text string)) *Textfield {
	field := new(Textfield)
	field.Name = name
	field.Text = text
	field.onChange = onChange
	field.currentContent = ""
	field.GuiType = "gui.TextField"
	return field
}

func (text *Textfield) Type() string {
	return "gui.TextField"
}

// GetContent returns the current content what was typed by the *Gui user.
func (text *Textfield) GetContent() string {
	return text.currentContent
}

// handleInputEvent handles the input event from the *Gui.
func (text *Textfield) handleInputEvent(gui *Gui) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := gui.AuthorizeOrRedirect(w, r)
		if err != nil {
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var inputData textInputData
		if err = json.Unmarshal(data, &inputData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		text.currentContent = inputData.Text
		text.onChange(username, text.currentContent)
	}
}

// AddToGui registers all handlers to the *Gui
func (text *Textfield) AddToGui(mount string, gui *Gui) {
	text.UpdateRequest = mount + text.Name + "/textfield/input"
	_ = gui.addURLFunc(text.UpdateRequest, text.handleInputEvent(gui))
}

// RemoveFromGui removes the component from the gui
func (text *Textfield) RemoveFromGui(mount string, gui *Gui) {
	text.UpdateRequest = mount + text.Name + "/textfield/input"
	_ = gui.removeURLFunc(text.UpdateRequest)
}
