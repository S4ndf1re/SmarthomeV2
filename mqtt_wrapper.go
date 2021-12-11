package main

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
)

type MQTTWrapper struct {
	config MQTTConfig
	client mqtt.Client
	mutex  sync.Mutex
}

func NewMQTTWrapper(config MQTTConfig) *MQTTWrapper {
	wrapper := new(MQTTWrapper)
	wrapper.config = config
	wrapper.mutex = sync.Mutex{}
	wrapper.client = mqtt.NewClient(config.ToPahoClientOptions())
	if token := wrapper.client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error: %s\n", token.Error())
	}
	return wrapper
}

func (wrapper *MQTTWrapper) Subscribe(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	if len(call.Arguments) == 2 {
		function, isFunction := goja.AssertFunction(call.Argument(1))
		if !isFunction {
			return vm.ToValue(errors.New("argument 2 is not a function"))
		}
		var topic string
		if err := vm.ExportTo(call.Argument(0), &topic); err != nil {
			return vm.ToValue(err)
		}
		callback := func(client mqtt.Client, data mqtt.Message) {
			_, _ = function(call.This, vm.ToValue(data.Topic()), vm.ToValue(string(data.Payload())))
		}
		if token := wrapper.client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
			fmt.Printf("%s\n", token.Error())
		}
	}
	return vm.ToValue(errors.New(fmt.Sprintf("expected 2 arguments. Got %d instead", len(call.Arguments))))
}

func (wrapper *MQTTWrapper) Publish(topic string, payload string, retained bool) {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	if token := wrapper.client.Publish(topic, 0, retained, payload); token.Wait() && token.Error() != nil {
		fmt.Printf("%s\n", token.Error())
	}
}

func (wrapper *MQTTWrapper) Unsubscribe(topic string) {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	wrapper.client.Unsubscribe(topic)
}

func (wrapper *MQTTWrapper) Close() {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	wrapper.client.Disconnect(1000)
}
