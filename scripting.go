package main

import (
	"github.com/dop251/goja"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ScriptDescriptor struct {
	Filepath   string `json:"filepath"`
	Sourcecode string `json:"sourcecode"`
	vm         *goja.Runtime
}

func LoadAllScripts(root string) []ScriptDescriptor {
	descriptors := make([]ScriptDescriptor, 0)

	_ = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && info.Name()[len(info.Name())-3:] == ".js" {
			if script, err := NewScript(path); err == nil {
				descriptors = append(descriptors, script)
			}
		}
		return nil
	})

	return descriptors
}

func NewScript(path string) (ScriptDescriptor, error) {
	descriptor := ScriptDescriptor{}
	descriptor.Filepath = path

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

func (script *ScriptDescriptor) run() error {
	_, err := script.vm.RunScript(script.Filepath, script.Sourcecode)
	return err
}

func (script *ScriptDescriptor) registerObjects() {
	_ = script.vm.Set("MQTTConfig", script.MQTTConfig)
	_ = script.vm.Set("MQTTWrapper", script.MQTTWrapper)

}

func (script *ScriptDescriptor) registerFunctions() {
	_ = script.vm.Set("GetScriptInstance", script.GetScriptInstance)
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
