package scripting

import (
	"fmt"
	"github.com/dop251/goja"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sync"
)

// MQTTWrapper is the wrapper around a Paho MQTT client
type MQTTWrapper struct {
	config MQTTConfig
	client mqtt.Client
	vm     *goja.Runtime
	mutex  sync.Mutex
}

// NewMQTTWrapper creates and connects a new *MQTTWrapper
func NewMQTTWrapper(config MQTTConfig, vm *goja.Runtime) *MQTTWrapper {
	wrapper := new(MQTTWrapper)
	wrapper.config = config
	wrapper.vm = vm
	wrapper.mutex = sync.Mutex{}
	wrapper.client = mqtt.NewClient(config.ToPahoClientOptions())
	if token := wrapper.client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("Error: %s\n", token.Error())
	}
	return wrapper
}

// Subscribe subscribes to a mqtt topic. If a publish is received for the subscribed topic, the callback accepted as the first argument is executed
// This function should only get called from the javascript code
// Javascript Example:
//
//	client.Subscribe("topic", (topic, payload) => { Print(payload) }
func (wrapper *MQTTWrapper) Subscribe(topic string, callback func(string, string)) error {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	mqttCallback := func(client mqtt.Client, data mqtt.Message) {
		callback(data.Topic(), string(data.Payload()))
	}
	if token := wrapper.client.Subscribe(topic, 0, mqttCallback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Publish publishes the payload to a topic. If retained is true, the mqtt message is marked as retained.
// This function is a javascript function it may be used from go code, but is supposed to get called from javascript
func (wrapper *MQTTWrapper) Publish(topic string, payload string, retained bool) error {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	if token := wrapper.client.Publish(topic, 0, retained, payload); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Unsubscribe unsubscribes from topic. The callback for that topic will get unregistered
func (wrapper *MQTTWrapper) Unsubscribe(topic string) {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	wrapper.client.Unsubscribe(topic)
}

// Close will end the Paho Client connection to the server
func (wrapper *MQTTWrapper) Close() {
	wrapper.mutex.Lock()
	defer wrapper.mutex.Unlock()
	wrapper.client.Disconnect(1000)
}
