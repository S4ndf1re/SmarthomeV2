package gui

import (
	"log"
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
	button.GuiType = "gui.Button"
	return button
}

// Type returns the type "gui.Button"
func (btn *Button) Type() string {
	return btn.GuiType
}

// AddToGui registers the Button onClick callback to the *Gui
func (btn *Button) AddToGui(mount string, gui *Gui) {
	btn.OnClickRequest = mount + btn.Name + "/button/click"
	_ = gui.addURLFunc(btn.OnClickRequest, func(w http.ResponseWriter, r *http.Request) {
		username, err := gui.AuthorizeOrRedirect(w, r)
		if err != nil {
			log.Printf("Button.AddToGui: %s\n", err)
			return
		}
		btn.onClick(username)
	})
}

// RemoveFromGui removes all handlers from the *Gui
func (btn *Button) RemoveFromGui(mount string, gui *Gui) {
	removePath := mount + btn.Name + "/button/click"
	if err := gui.removeURLFunc(removePath); err != nil {
		log.Printf("Button.RemoveFromGui(): %s\n", err)
	}
}
