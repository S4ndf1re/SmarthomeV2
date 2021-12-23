package gui

import (
	"Smarthome/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type textInputData struct {
	Text string `json:"text"`
}

// TextField is a text input component
type TextField struct {
	Name           string `json:"name"`
	Text           string `json:"text"`
	UpdateRequest  string `json:"updateRequest"`
	GuiType        string `json:"type"`
	onChange       func(username string, text string)
	currentContent string
}

// NewTextField constructs a new text field with a default change listener
// Internal state is kept track of
func NewTextField(name string, text string, onChange func(username string, text string)) *TextField {
	field := new(TextField)
	field.Name = name
	field.Text = text
	field.onChange = onChange
	field.currentContent = ""
	field.GuiType = TextFieldType
	return field
}

func (text *TextField) Type() string {
	return text.GuiType
}

// GetContent returns the current content what was typed by the *Gui user.
func (text *TextField) GetContent() string {
	return text.currentContent
}

// handleInputEvent handles the input event from the *Gui.
func (text *TextField) handleInputEvent(username string, w http.ResponseWriter, r *http.Request) {
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

// AddToGui registers all handlers to the *Gui
func (text *TextField) AddToGui(mount string, gui *Gui) {
	text.UpdateRequest = mount + text.Name + textFieldTextInputRequest
	err := gui.addURLFunc(text.UpdateRequest, gui.AuthorizeOrRedirect(text.handleInputEvent))
	util.LogIfErr("TextField.AddToGui()", err)
}

// RemoveFromGui removes the component from the gui
func (text *TextField) RemoveFromGui(mount string, gui *Gui) {
	text.UpdateRequest = mount + text.Name + textFieldTextInputRequest
	err := gui.removeURLFunc(text.UpdateRequest)
	util.LogIfErr("TextField.RemoveFromGui()", err)
}
