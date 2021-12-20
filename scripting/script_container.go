package scripting

import (
	"Smarthome/gui"
	"sync"
)

// ScriptContainer contains all scripts
type ScriptContainer struct {
	root    string
	gui     *gui.Gui
	scripts []*ScriptDescriptor
}

// NewScriptContainer creates a new ScriptContainer
// All scripts are loaded from the root directory
func NewScriptContainer(root string, gui *gui.Gui) ScriptContainer {
	return ScriptContainer{
		scripts: LoadAllScripts(root, gui),
		gui:     gui,
	}
}

// RunAll executes all scripts in a goroutine for each script
func (container *ScriptContainer) RunAll() *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	for _, script := range container.scripts {
		go script.run(wg)
	}
	return wg
}

// ReloadAllScripts reloads all scripts. Before reloading, the close function is from the javascript script is called
func (container *ScriptContainer) ReloadAllScrips() {

	newScripts := make([]*ScriptDescriptor, 0)

	for _, oldScript := range container.scripts {
		oldScript.Close()
		if newScript, err := NewScript(oldScript.Filepath, oldScript.Name, container.gui); err == nil {
			newScripts = append(newScripts, newScript)
		}
	}
	container.scripts = newScripts
}
