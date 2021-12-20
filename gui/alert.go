package gui

// Alert is a simple alert message
type Alert struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	GuiType string `json:"type"`
}

// NewAlert generates a new alert. If message is empty, the Alert will not get triggered in the gui
func NewAlert(name, message string) *Alert {
	alert := new(Alert)
	alert.Name = name
	alert.Message = message
	alert.GuiType = "gui.Alert"

	return alert
}

func (alert *Alert) Type() string {
	return "gui.Alert"
}

// AddToGui stub
func (alert *Alert) AddToGui(_ string, _ *Gui) {

}

// RemoveFromGui stub
func (alert *Alert) RemoveFromGui(_ string, _ *Gui) {

}
