package gui

// Child represents a simple Gui child. Every child must be Addable
type Child interface {
	Type() string
	Addable
}
