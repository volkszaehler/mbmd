package server

import (
	"fmt"
	"github.com/volkszaehler/mbmd/prometheus"
	"log"
	"regexp"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/volkszaehler/mbmd/meters"
)

const (
	publishTimeout = 2000 * time.Millisecond
)

var (
	topicRE = regexp.MustCompile(`(\w+)([LTS]\d)`)
)

// MqttClient is a MQTT publisher
type MqttClient struct {
	Client  MQTT.Client
	qos     byte
	verbose bool
}

// NewMqttOptions creates MQTT client options
func NewMqttOptions(
	broker string,
	user string,
	password string,
	clientID string,
) *MQTT.ClientOptions {
	opt := MQTT.NewClientOptions()
	opt.AddBroker(broker)
	opt.SetUsername(user)
	opt.SetPassword(password)
	opt.SetClientID(clientID)
	opt.SetAutoReconnect(true)
	return opt
}

// NewMqttClient creates new publisher for MQTT
func NewMqttClient(
	options *MQTT.ClientOptions,
	qos byte,
	verbose bool,
) *MqttClient {
	log.Printf("mqtt: connecting %s at %s", options.ClientID, options.Servers)

	client := MQTT.NewClient(options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		prometheus.PublisherConnectionFailure.WithLabelValues("mqtt").Inc()
		log.Fatalf("mqtt: error connecting: %s", token.Error())
	} else if token.Wait() && token.Error() == nil {
		prometheus.PublisherConnectionSuccess.WithLabelValues("mqtt").Inc()
	}

	if verbose {
		log.Println("mqtt: connected")
	}

	prometheus.PublisherCreated.WithLabelValues("mqtt").Inc()

	return &MqttClient{
		Client:  client,
		qos:     qos,
		verbose: verbose,
	}
}

// Publish MQTT message with error handling
func (m *MqttClient) Publish(topic string, retained bool, message interface{}) {
	token := m.Client.Publish(topic, m.qos, retained, message)
	prometheus.PublisherDataPublishAttempt.WithLabelValues("mqtt").Inc()
	if m.verbose {
		log.Printf("mqtt: publish %s, message: %s", topic, message)
	}
	go m.WaitForToken(token)
}

// WaitForToken synchronously waits until token operation completed
func (m *MqttClient) WaitForToken(token MQTT.Token) {
	if token.WaitTimeout(publishTimeout) {
		if token.Error() != nil {
			log.Printf("mqtt: error: %s", token.Error())
			prometheus.PublisherDataPublishedError.WithLabelValues("mqtt").Inc()
		} else {
			prometheus.PublisherDataPublished.WithLabelValues("mqtt").Inc()
		}
	} else {
		prometheus.PublisherConnectionTimeOut.WithLabelValues("mqtt").Inc()
		if m.verbose {
			log.Println("mqtt: timeout")
		}
	}
}

// deviceTopic converts meter's device id to topic string
func mqttDeviceTopic(deviceID string) string {
	topic := strings.Replace(strings.ToLower(deviceID), "#", "", -1)
	return strings.Replace(topic, ".", "-", -1)
}

// MqttRunner allows to attach an MqttClient as broadcast receiver
type MqttRunner struct {
	*MqttClient
	topic string
}

// NewMqttRunner create a new runer for plain MQTT
func NewMqttRunner(options *MQTT.ClientOptions, qos byte, topic string, verbose bool) *MqttRunner {
	// set will
	lwt := fmt.Sprintf("%s/status", topic)
	options.SetWill(lwt, "disconnected", qos, true)

	client := NewMqttClient(options, qos, verbose)

	return &MqttRunner{
		MqttClient: client,
		topic:      topic,
	}
}

// topicFromMeasurement converts measurements of type MeasureLx/MeasureSx/MeasureTx to hierarchical Measure/Lx topics
func topicFromMeasurement(measurement meters.Measurement) string {
	name := measurement.String()
	match := topicRE.FindStringSubmatch(name)
	if len(match) != 3 {
		return name
	}

	topic := fmt.Sprintf("%s/%s", match[1], match[2])

	return topic
}

// Run MqttClient publisher
func (m *MqttRunner) Run(in <-chan QuerySnip) {
	// notify connection and override will
	m.MqttClient.Publish(fmt.Sprintf("%s/status", m.topic), true, "connected")

	for snip := range in {
		subtopic := topicFromMeasurement(snip.Measurement)
		topic := fmt.Sprintf("%s/%s/%s", m.topic, mqttDeviceTopic(snip.Device), subtopic)
		message := fmt.Sprintf("%.3f", snip.Value)
		go m.Publish(topic, false, message)
	}
}
