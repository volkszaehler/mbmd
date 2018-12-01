package sdm630

import (
	"fmt"
	"log"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type RateMap map[string]int64

// Allowed checks if topic has bee published longer than rate ago
func (r *RateMap) Allowed(rate int, topic string) bool {
	if rate == 0 {
		return true
	}

	t := (*r)[topic]
	now := time.Now().Unix()
	if now > t {
		(*r)[topic] = now + int64(rate)
		return true
	}

	return false
}

type MqttClient struct {
	client    MQTT.Client
	mqttTopic string
	mqttQos   int
	verbose   bool
}

// Publish MQTT message with error handling
func (m *MqttClient) Publish(topic string, retained bool, message interface{}) {
	token := m.client.Publish(topic, byte(m.mqttQos), false, message)
	if m.verbose {
		log.Printf("MQTT: publish %s, message: %s", topic, message)
	}

	if token.WaitTimeout(2000 * time.Millisecond) {
		if token.Error() != nil {
			log.Printf("MQTT: error: %s", token.Error())
		}
	} else {
		if m.verbose {
			log.Printf("MQTT: timeout")
		}
	}
}

func NewMqttClient(
	mqttBroker string,
	mqttTopic string,
	mqttUser string,
	mqttPassword string,
	mqttClientID string,
	mqttQos int,
	mqttCleanSession bool,
	verbose bool,
) *MqttClient {
	mqttOpts := MQTT.NewClientOptions()
	mqttOpts.AddBroker(mqttBroker)
	mqttOpts.SetUsername(mqttUser)
	mqttOpts.SetPassword(mqttPassword)
	mqttOpts.SetClientID(mqttClientID)
	mqttOpts.SetCleanSession(mqttCleanSession)
	mqttOpts.SetAutoReconnect(true)

	topic := fmt.Sprintf("%s/status", mqttTopic)
	message := fmt.Sprintf("disconnected")
	mqttOpts.SetWill(topic, message, byte(mqttQos), true)

	log.Printf("Connecting MQTT at %s", mqttBroker)
	if verbose {
		log.Printf("\tclientid:     %s\n", mqttClientID)
		log.Printf("\tuser:         %s\n", mqttUser)
		if mqttPassword != "" {
			log.Printf("\tpassword:     ****\n")
		}
		log.Printf("\ttopic:        %s\n", mqttTopic)
		log.Printf("\tqos:          %d\n", mqttQos)
		log.Printf("\tcleansession: %v\n", mqttCleanSession)
	}

	mqttClient := MQTT.NewClient(mqttOpts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("MQTT: error connecting: ", token.Error())
		panic(token.Error())
	}
	if verbose {
		log.Println("MQTT: connected")
	}

	// notify connection
	message = fmt.Sprintf("connected")
	token := mqttClient.Publish(topic, byte(mqttQos), true, message)
	if verbose {
		log.Printf("MQTT: publish %s, message: %s", topic, message)
	}
	if token.Wait() && token.Error() != nil {
		log.Fatal("MQTT: error connecting, trying to reconnect: ", token.Error())
	}

	return &MqttClient{
		client:    mqttClient,
		mqttTopic: mqttTopic,
		mqttQos:   mqttQos,
		verbose:   verbose,
	}
}

// DeviceTopic converts meter's device id to topic string
func (m *MqttClient) DeviceTopic(deviceId uint8) string {
	uniqueID := fmt.Sprintf(UniqueIdFormat, deviceId)
	return strings.Replace(strings.ToLower(uniqueID), "#", "", -1)
}

type MqttRunner struct {
	*MqttClient
}

// Run MQTT client publisher
func (m *MqttRunner) Run(in QuerySnipChannel, rate int) {
	rateMap := make(RateMap)

	for {
		snip := <-in
		topic := fmt.Sprintf("%s/%s/%s", m.mqttTopic, m.DeviceTopic(snip.DeviceId), snip.IEC61850)

		if rateMap.Allowed(rate, topic) {
			message := fmt.Sprintf("%.3f", snip.Value)
			go m.Publish(topic, false, message)
		} else {
			if m.verbose {
				log.Printf("MQTT: skipped %s, rate to high", topic)
			}
		}
	}
}
