package gui

// Addable can be added to a gui as a http handler
type Addable interface {
	AddToGui(mount string, gui *Gui)
	RemoveFromGui(mount string, gui *Gui)
}
