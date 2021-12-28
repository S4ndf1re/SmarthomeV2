package gui

// Alert is a simple alert message
type Alert struct {
	Name     string `json:"name"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
	GuiType  string `json:"type"`
}

// NewAlert generates a new alert. If message is empty, the Alert will not get triggered in the gui
func NewAlert(name, message, severity string) *Alert {
	alert := new(Alert)
	alert.Name = name
	alert.Message = message
	alert.Severity = severity
	alert.GuiType = AlertType

	return alert
}

func (alert *Alert) Type() string {
	return alert.GuiType
}

// AddToGui stub
func (alert *Alert) AddToGui(_ string, _ *Gui) {

}

// RemoveFromGui stub
func (alert *Alert) RemoveFromGui(_ string, _ *Gui) {

}
