package main

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTWrapper struct {
	config MQTTConfig
	client mqtt.Client
}

func NewMQTTWrapper(config MQTTConfig) *MQTTWrapper {
	wrapper := new(MQTTWrapper)
	wrapper.config = config
	wrapper.client = mqtt.NewClient(config.ToPahoClientOptions())
	return wrapper
}

func (wrapper *MQTTWrapper) Subscribe(call goja.FunctionCall, vm *goja.Runtime) goja.Value {
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
			_, _ = function(call.This, vm.ToValue(data.Topic()), vm.ToValue(data.Payload()))
		}
		wrapper.client.Subscribe(topic, 0, callback)
	}
	return vm.ToValue(errors.New(fmt.Sprintf("expected 2 arguments. Got %d instead", len(call.Arguments))))
}

func (wrapper *MQTTWrapper) Publish(topic string, payload []byte, retained bool) {
	wrapper.client.Publish(topic, 0, retained, payload)
}

func (wrapper *MQTTWrapper) Unsubscribe(topic string) {
	wrapper.client.Unsubscribe(topic)
}

func (wrapper *MQTTWrapper) Close() {
	wrapper.client.Disconnect(1000)
}
