package api

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/mainflux/mainflux/logger"
)

type MQTTPublisher interface {
	connect() error
	Publish([]byte)
}

type mqttPublisher struct {
	logger     *log.Logger
	mqttClient mqtt.Client
	host       string
	username   string
	password   string
	clientID   string
	topic      string
}

var _ MQTTPublisher = (*mqttPublisher)(nil)

// NewMQTTPublisher creates a new MQTT publisher, configures it's topic and connects it to the MQTT server.
func NewMQTTPublisher(logger *log.Logger, host, username, password, clientID, topic string) (*mqttPublisher, error) {
	p := mqttPublisher{logger: logger, host: host, username: username, password: password, clientID: clientID, topic: topic}
	if err := p.connect(); err != nil {
		return &p, err
	}
	return &p, nil
}

func (m *mqttPublisher) connect() error {
	opts := mqtt.NewClientOptions().AddBroker(m.host)
	opts.SetUsername(m.username)
	opts.SetPassword(m.password)
	opts.SetClientID(m.clientID)

	mc := mqtt.NewClient(opts)

	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	m.mqttClient = mc
	return nil
}

// Publish publishes a message to the configured topic.
func (m *mqttPublisher) Publish(payload []byte) {

	(*m.logger).Debug(fmt.Sprintf("Publishing message %s", payload))
	if token := m.mqttClient.Publish(m.topic, 0, false, payload); token.Wait() && token.Error() != nil {
		(*m.logger).Error(fmt.Sprintf("Failed to publish message on topic %s : %s", m.topic, token.Error()))
	}

}
