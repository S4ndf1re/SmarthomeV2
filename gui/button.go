package gui

import (
	"Smarthome/util"
	"net/http"
)

// Button is a simple gui button with onClick handler
type Button struct {
	Name           string `json:"name"`
	Text           string `json:"text"`
	OnClickRequest string `json:"onClickRequest"`
	onClick        func(user string)
	GuiType        string `json:"type"`
}

// NewButton creates a new button and registers the [onClick] handler
func NewButton(name, text string, onClick func(user string)) *Button {
	button := new(Button)
	button.Name = name
	button.Text = text
	button.OnClickRequest = ""
	button.onClick = onClick
	button.GuiType = ButtonType
	return button
}

// Type returns the type "gui.Button"
func (btn *Button) Type() string {
	return btn.GuiType
}

func (btn *Button) handleOnClick(username string, _ http.ResponseWriter, _ *http.Request) {
	btn.onClick(username)
}

// AddToGui registers the Button onClick callback to the *Gui
func (btn *Button) AddToGui(mount string, gui *Gui) {
	btn.OnClickRequest = mount + btn.Name + buttonOnClickExtension
	err := gui.addURLFunc(btn.OnClickRequest, gui.AuthorizeOrRedirect(btn.handleOnClick))
	util.LogIfErr("Button.AddToGui()", err)
}

// RemoveFromGui removes all handlers from the *Gui
func (btn *Button) RemoveFromGui(mount string, gui *Gui) {
	btn.OnClickRequest = mount + btn.Name + buttonOnClickExtension
	err := gui.removeURLFunc(btn.OnClickRequest)
	util.LogIfErr("Button.RemoveFromGui()", err)
}
