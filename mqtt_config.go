package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"strconv"
)

type MQTTConfig struct {
	LastWill        bool   `json:"lastWill"`
	LastWillMessage []byte `json:"lastWillMessage"`
	LastWillRetain  bool   `json:"lastWillRetain"`
	hostname        string
	port            int
}

func (config MQTTConfig) ToPahoClientOptions() *mqtt.ClientOptions {
	options := mqtt.NewClientOptions()
	options.Servers = []*url.URL{
		{
			Host: config.hostname + ":" + strconv.Itoa(config.port),
		},
	}
	options.WillEnabled = config.LastWill
	options.WillPayload = config.LastWillMessage
	options.WillRetained = config.LastWillRetain

	return options
}
