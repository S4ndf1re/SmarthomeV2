package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConfig struct {
	LastWill        bool   `json:"lastWill"`
	LastWillTopic   string `json:"lastWillTopic"`
	LastWillMessage string `json:"lastWillMessage"`
	LastWillRetain  bool   `json:"lastWillRetain"`
	Hostname        string `json:"hostname"`
	Port            int    `json:"port"`
}

func (config MQTTConfig) ToPahoClientOptions() *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Hostname, config.Port))
	if config.LastWill {
		options.SetWill(config.LastWillTopic, config.LastWillMessage, 0, config.LastWillRetain)
	}

	return options
}
