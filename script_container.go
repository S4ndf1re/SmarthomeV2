package main

import "sync"

type ScriptContainer struct {
	root    string
	scripts []*ScriptDescriptor
}

func NewScriptContainer(root string) ScriptContainer {
	return ScriptContainer{
		scripts: LoadAllScripts(root),
	}
}

func (container *ScriptContainer) RunAll() *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	for _, script := range container.scripts {
		go script.run(wg)
	}
	return wg
}

func (container *ScriptContainer) ReloadAllScrips() {

	newScripts := make([]*ScriptDescriptor, 0)

	for _, oldScript := range container.scripts {
		oldScript.Close()
		if newScript, err := NewScript(oldScript.Filepath, oldScript.Name); err == nil {
			newScripts = append(newScripts, newScript)
		}
	}
	container.scripts = newScripts
}
