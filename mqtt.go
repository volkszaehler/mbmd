package sdm630

import (
	"fmt"
	"log"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client    MQTT.Client
	mqttTopic string
	mqttRate  int
	mqttQos   int
	in        QuerySnipChannel
	verbose   bool
}

// Run MQTT client publisher
func (m *MqttClient) Run() {
	mqttRateMap := make(map[string]int64)

	for {
		snip := <-m.in
		topic := fmt.Sprintf("%s/%s/%s", m.mqttTopic, m.MeterTopic(snip.DeviceId), snip.IEC61850)

		t := mqttRateMap[topic]
		now := time.Now()
		if m.mqttRate == 0 || now.Unix() > t {
			message := fmt.Sprintf("%.3f", snip.Value)
			go m.Publish(topic, false, message)
			mqttRateMap[topic] = now.Unix() + int64(m.mqttRate)
		} else {
			if m.verbose {
				log.Printf("MQTT: skipped %s, rate to high", topic)
			}
		}
	}
}

// MeterTopic converts meter's device id to topic string
func (m *MqttClient) MeterTopic(deviceId uint8) string {
	uniqueID := fmt.Sprintf(UniqueIdFormat, deviceId)
	return strings.Replace(strings.ToLower(uniqueID), "#", "", -1)
}

// Publish MQTT message with error handling
func (m *MqttClient) Publish(topic string, retained bool, message interface{}) {
	token := m.client.Publish(topic, byte(m.mqttQos), false, message)
	if m.verbose {
		log.Printf("MQTT: publish %s, message: %s", topic, message)
	}

	if token.WaitTimeout(2000 * time.Millisecond) {
		if token.Error() != nil {
			log.Printf("MQTT: Error: %s", token.Error())
		}
	} else {
		if m.verbose {
			log.Printf("MQTT: Timeout")
		}
	}
}

func NewMqttClient(
	in QuerySnipChannel,
	mqttBroker string,
	mqttTopic string,
	mqttUser string,
	mqttPassword string,
	mqttClientID string,
	mqttQos int,
	mqttRate int,
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
		log.Fatal("MQTT: Error connecting, trying to reconnect: ", token.Error())
	}

	return &MqttClient{
		in:        in,
		client:    mqttClient,
		mqttTopic: mqttTopic,
		mqttRate:  mqttRate,
		mqttQos:   mqttQos,
		verbose:   verbose,
	}
}
