package main

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func configureMqttConnection() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%s", MQTT_URL, MQTT_PORT))
	opts.SetClientID("back-end")
	opts.SetUsername(MQTT_USERNAME)
	opts.SetPassword(MQTT_PASSWORD)
	return opts
}

func NewMqttClient() (mqtt.Client, error) {
	opts := configureMqttConnection()
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

var payload Data

func (app *Config) Subscribe(channel chan Data) error {
	if token := app.MqttClient.Subscribe(MQTT_TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {
		if json.Unmarshal(msg.Payload(), &payload) == nil {
			channel <- payload
		}
	}); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (app *Config) Unsubscribe() {
	app.MqttClient.Unsubscribe(MQTT_TOPIC)
}

func (app *Config) Disconnect(millis uint) {
	app.MqttClient.Disconnect(millis)
}
