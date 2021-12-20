package scripting

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTConfig ist the configuration for a MQTTWrapper
// Not all MQTTWrapper options are accessible. This is because the javascript should not be allowed to configure every single option
type MQTTConfig struct {
	LastWill        bool   `json:"lastWill"`
	LastWillTopic   string `json:"lastWillTopic"`
	LastWillMessage string `json:"lastWillMessage"`
	LastWillRetain  bool   `json:"lastWillRetain"`
	Hostname        string `json:"hostname"`
	Port            int    `json:"port"`
}

// ToPahoClientOptions converts the MQTTConfig to Paho MQTT client options
func (config MQTTConfig) ToPahoClientOptions() *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Hostname, config.Port))
	if config.LastWill {
		options.SetWill(config.LastWillTopic, config.LastWillMessage, 0, config.LastWillRetain)
	}

	return options
}
