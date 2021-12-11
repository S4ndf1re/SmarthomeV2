package main

import (
	"github.com/dop251/goja"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ScriptDescriptor struct {
	Filepath   string `json:"filepath"`
	Name       string `json:"name"`
	Sourcecode string `json:"sourcecode"`
	mutex      sync.Mutex
	vm         *goja.Runtime
}

func LoadAllScripts(root string) []*ScriptDescriptor {
	descriptors := make([]*ScriptDescriptor, 0)

	_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && info.Name()[len(info.Name())-3:] == ".js" {
			if script, err := NewScript(path, info.Name()[:len(info.Name())-3]); err == nil {
				descriptors = append(descriptors, script)
			}
		}
		return nil
	})

	return descriptors
}

func NewScript(path, name string) (*ScriptDescriptor, error) {
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

	return descriptor, nil
}

func (script *ScriptDescriptor) run(wg *sync.WaitGroup) {
	script.mutex.Lock()
	wg.Add(1)
	defer wg.Done()
	defer script.mutex.Unlock()
	_, _ = script.vm.RunScript(script.Filepath, script.Sourcecode)
}

func (script *ScriptDescriptor) registerObjects() {
	_ = script.vm.Set("MQTTConfig", script.MQTTConfig)
	_ = script.vm.Set("MQTTWrapper", script.MQTTWrapper)

}

func (script *ScriptDescriptor) registerFunctions() {
	_ = script.vm.Set("GetScriptInstance", script.GetScriptInstance)
	_ = script.vm.Set("ReadFile", script.ReadFile)
	_ = script.vm.Set("WriteFile", script.WriteFile)
	RegisterToGojaVM(script.vm)
}

func (script *ScriptDescriptor) GetScriptInstance() *goja.Object {
	return script.vm.ToValue(script).(*goja.Object)
}

func (script *ScriptDescriptor) MQTTConfig(call goja.ConstructorCall) *goja.Object {
	instance := script.vm.ToValue(new(MQTTConfig)).(*goja.Object)
	_ = instance.SetPrototype(call.This.Prototype())
	return instance
}

func (script *ScriptDescriptor) MQTTWrapper(call goja.ConstructorCall) *goja.Object {
	if len(call.Arguments) == 1 {
		config := new(MQTTConfig)
		_ = script.vm.ExportTo(call.Argument(0), config)
		wrapper := NewMQTTWrapper(*config)
		instance := script.vm.ToValue(wrapper).(*goja.Object)
		_ = instance.SetPrototype(call.This.Prototype())
		return instance
	}
	return nil
}

func (script *ScriptDescriptor) ReadFile(path string) (string, error) {
	strings.ReplaceAll(path, "..", ".")
	path = "scriptfiles" + string(filepath.Separator) + script.Name + string(filepath.Separator) + path
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

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
	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		return err
	}

	return nil
}

func (script *ScriptDescriptor) Close() {
	var close func()
	script.vm.ExportTo(script.vm.Get("close"), &close)
	close()
	script.vm.Interrupt("Dispose engine")
}
