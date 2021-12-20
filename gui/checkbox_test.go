package gui

func ExampleCheckbox_OverrideListeners() {
	state := false
	checkbox := NewCheckbox("test", "Check Me")

	checkbox.OverrideListeners(func(s string) {
		// OnOffState
		state = false
	}, func(s string) {
		// OnOnState
		state = true
	}, func(s string) bool {
		// OnGetState
		return state
	})

	// Register defaults again:
	checkbox.OverrideListeners(checkbox.DefaultListeners())
}
