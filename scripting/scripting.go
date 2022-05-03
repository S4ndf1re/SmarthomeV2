package scripting

import (
	"Smarthome/gui"
	"Smarthome/util"
	"github.com/dop251/goja"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ScriptDescriptor contains all important script information
// Each javascript script has access to the ScriptDescriptor by calling the GetScriptInstance function
type ScriptDescriptor struct {
	Filepath   string `json:"filepath"`
	Name       string `json:"name"`
	Sourcecode string `json:"sourcecode"`
	gui        *gui.Gui
	mutex      sync.Mutex
	vm         *goja.Runtime
}

// LoadAllScripts loads all scripts that are contained in the root folder
func LoadAllScripts(root string, gui *gui.Gui) []*ScriptDescriptor {
	descriptors := make([]*ScriptDescriptor, 0)

	_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && info.Name()[len(info.Name())-3:] == ".js" {
			if script, err := NewScript(path, info.Name()[:len(info.Name())-3], gui); err == nil {
				descriptors = append(descriptors, script)
			}
		}
		return nil
	})

	return descriptors
}

// NewScript initializes a new script. This includes a new *goja.Runtime. Also, all functions and objects are getting registered
func NewScript(path, name string, gui *gui.Gui) (*ScriptDescriptor, error) {
	descriptor := new(ScriptDescriptor)
	descriptor.Filepath = path
	descriptor.Name = name
	descriptor.mutex = sync.Mutex{}

	vm := goja.New()

	file, err := os.Open(path)
	if err != nil {
		return descriptor, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return descriptor, err
	}

	descriptor.Sourcecode = string(data)
	descriptor.vm = vm

	descriptor.registerObjects()
	descriptor.registerFunctions()
	descriptor.gui = gui

	return descriptor, nil
}

// run executes the script contained in the ScriptDescriptor
func (script *ScriptDescriptor) run(wg *sync.WaitGroup) {
	script.mutex.Lock()
	wg.Add(1)
	defer wg.Done()
	defer script.mutex.Unlock()
	log.Printf("Running script %s\n", script.Name)
	if _, err := script.vm.RunScript(script.Filepath, script.Sourcecode); err != nil {
		log.Printf("vm.RunScript(): %s\n", err)
	} else {
		log.Printf("Done running script %s\n", script.Name)
	}
}

// registerObjects registers all needed and public Objects to the goja.Runtime of the ScriptDescriptor
func (script *ScriptDescriptor) registerObjects() {
	_ = script.vm.Set("MQTTConfig", script.MQTTConfig)
	_ = script.vm.Set("MQTTWrapper", script.MQTTWrapper)
	_ = script.vm.Set("Container", script.Container)
	_ = script.vm.Set("Button", script.Button)
	_ = script.vm.Set("Checkbox", script.Checkbox)
	_ = script.vm.Set("Alert", script.Alert)
	_ = script.vm.Set("TextField", script.TextField)
	_ = script.vm.Set("Data", script.Data)
	_ = script.vm.Set("TCPClient", script.TCPClient)
}

// registerFunctions registers all needed and public functions to the goja.Runtime of the ScriptDescriptor
func (script *ScriptDescriptor) registerFunctions() {
	_ = script.vm.Set("GetScriptInstance", script.GetScriptInstance)
	_ = script.vm.Set("ReadFile", script.ReadFile)
	_ = script.vm.Set("WriteFile", script.WriteFile)
	_ = script.vm.Set("AddContainer", script.AddContainer)
	_ = script.vm.Set("RemoveContainer", script.RemoveContainer)
	util.RegisterToGojaVM(script.vm)
}

// GetScriptInstance is a function that can be called from the javascript script. It will return the corresponding
// ScriptDescriptor that belongs to the calling script
func (script *ScriptDescriptor) GetScriptInstance() *goja.Object {
	return script.vm.ToValue(script).(*goja.Object)
}

// MQTTConfig constructs (constructor call) a new MQTTConfig as a *goja.Object
func (script *ScriptDescriptor) MQTTConfig(call goja.ConstructorCall) *goja.Object {
	instance := script.vm.ToValue(new(MQTTConfig)).(*goja.Object)
	_ = instance.SetPrototype(call.This.Prototype())
	return instance
}

// MQTTWrapper constructs (constructor call) a new MQTTWrapper as a *goja.Object.
// The constructor call requires to have exactly one argument, otherwise nil is returned
// The single argument must be of type MQTTConfig
func (script *ScriptDescriptor) MQTTWrapper(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 1 {
		config := new(MQTTConfig)
		_ = script.vm.ExportTo(call.Argument(0), config)
		wrapper := NewMQTTWrapper(*config, script.vm)
		instance := script.vm.ToValue(wrapper).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

// ReadFile reads the file path from the scripts corresponding and isolated folder.
// .. is filtered out of the filepath.
// If the file does not exist, an error is returned
func (script *ScriptDescriptor) ReadFile(path string) (string, error) {
	strings.ReplaceAll(path, "..", ".")
	path = "scriptfiles" + string(filepath.Separator) + script.Name + string(filepath.Separator) + path
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("file.Close(): %s\n", err)
		}
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// WriteFile writes data to the file at path.
// The file is contained in the scripts corresponding folder.
// .. is filtered out from the filepath.
// If the file doesn't exist, the file is created. If it exists, it gets overwritten
func (script *ScriptDescriptor) WriteFile(path, data string) error {
	strings.ReplaceAll(path, "..", ".")
	newPath := filepath.Join("scriptfiles", script.Name)
	if err := os.MkdirAll(newPath, os.ModePerm); err != nil {
		return err
	}
	newPath = filepath.Join(newPath, path)
	file, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("file.Close(): %s\n", err)
		}
	}()

	if _, err = file.WriteString(data); err != nil {
		return err
	}

	return nil
}

// AddContainer adds a container to the internal *Gui. This function is exposed to javascript
func (script *ScriptDescriptor) AddContainer(container *gui.Container) {
	script.gui.AddContainer(container)
}

// RemoveContainer removes a container from the internal *Gui. This function is exposed to javascript
func (script *ScriptDescriptor) RemoveContainer(container *gui.Container) {
	script.gui.RemoveContainer(container)
}

// Close will call the close function of the javascript script if the function exists.
// After the close function is called, the vm is disposed
func (script *ScriptDescriptor) Close() {
	var closeFunction func()
	if err := script.vm.ExportTo(script.vm.Get("closeFunction"), &closeFunction); err != nil {
		closeFunction()
	}
	script.vm.Interrupt("Dispose engine")
}

// Container is the javascript constructor call for the gui.Container.
// Three constructor parameters are required: name string, text string and onInit func(string)
func (script *ScriptDescriptor) Container(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 4 {
		var name string
		var text string
		var onInit func(string)
		var onUnload func(string)
		if err := script.vm.ExportTo(call.Argument(0), &name); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(1), &text); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(2), &onInit); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(3), &onUnload); err != nil {
			return nil
		}
		container := gui.NewContainer(script.Name+"/"+name, text, onInit, onUnload)
		instance := script.vm.ToValue(container).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

// Button is the javascript constructor for gui.Button.
// Three arguments are required: name string, text string and onClick func(string)
func (script *ScriptDescriptor) Button(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 3 {
		var name string
		var text string
		var onClick func(string)
		if err := script.vm.ExportTo(call.Argument(0), &name); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(1), &text); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(2), &onClick); err != nil {
			return nil
		}
		button := gui.NewButton(name, text, onClick)
		instance := script.vm.ToValue(button).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

// Checkbox is the javascript constructor for gui.Checkbox.
// Two arguments are required: name string and text string
func (script *ScriptDescriptor) Checkbox(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 2 {
		var name string
		var text string
		if err := script.vm.ExportTo(call.Argument(0), &name); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(1), &text); err != nil {
			return nil
		}
		checkbox := gui.NewCheckbox(name, text)
		instance := script.vm.ToValue(checkbox).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

// Textfield is the javascript constructor for gui.TextField.
// Three arguments are required: name string, text string, onChange func(string, string)
func (script *ScriptDescriptor) TextField(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 3 {
		var name string
		var text string
		var onChange func(string, string)
		if err := script.vm.ExportTo(call.Argument(0), &name); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(1), &text); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(2), &onChange); err != nil {
			return nil
		}
		textField := gui.NewTextField(name, text, onChange)
		instance := script.vm.ToValue(textField).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

// Alert is the javascript constructor for gui.Alert.
// Two arguments are required: name string and message string
func (script *ScriptDescriptor) Alert(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 3 {
		var name string
		var message string
		var severity string
		if err := script.vm.ExportTo(call.Argument(0), &name); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(1), &message); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(2), &severity); err != nil {
			return nil
		}
		alert := gui.NewAlert(name, message, severity)
		instance := script.vm.ToValue(alert).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

// Data is the javascript constructor for gui.Data.
// Two arguments are required: name string and initialChild gui.Child
func (script *ScriptDescriptor) Data(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 2 {
		var name string
		var initialChild gui.Child
		if err := script.vm.ExportTo(call.Argument(0), &name); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(1), &initialChild); err != nil {
			return nil
		}
		data := gui.NewData(name, initialChild)
		instance := script.vm.ToValue(data).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

func (script *ScriptDescriptor) TCPClient(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 2 {
		var hostname string
		var port uint16
		if err := script.vm.ExportTo(call.Argument(0), &hostname); err != nil {
			return nil
		}
		if err := script.vm.ExportTo(call.Argument(1), &port); err != nil {
			return nil
		}
		data := NewTCPClient(hostname, port)
		instance := script.vm.ToValue(data).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

func (script *ScriptDescriptor) Mutex(call goja.ConstructorCall) *goja.Object {
	data := &sync.Mutex{}
	instance := script.vm.ToValue(data).(*goja.Object)
	_ = instance.SetPrototype(call.This.Prototype())
	return instance
}
