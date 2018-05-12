package sdm630

import (
	"fmt"
	"log"
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
func (mqttClient *MqttClient) Run() {
	mqttRateMap := make(map[string]int64)

	for {
		snip := <-mqttClient.in
		if mqttClient.verbose {
			log.Printf("MQTT: got meter data (device %d: data: %s, value: %.3f W, desc: %s)", snip.DeviceId, snip.IEC61850, snip.Value, snip.Description)
		}

		uniqueID := fmt.Sprintf(UniqueIdFormat, snip.DeviceId)
		topic := fmt.Sprintf("%s/%s/%s", mqttClient.mqttTopic, uniqueID, snip.IEC61850)

		t := mqttRateMap[topic]
		now := time.Now()
		if mqttClient.mqttRate == 0 || now.Unix() > t {
			message := fmt.Sprintf("%.3f", snip.Value)
			token := mqttClient.client.Publish(topic, byte(mqttClient.mqttQos), true, message)
			if mqttClient.verbose {
				log.Printf("MQTT: push %s, message: %s", topic, message)
			}
			if token.Wait() && token.Error() != nil {
				log.Fatal("MQTT: Error connecting, trying to reconnect: ", token.Error())
			}
			mqttRateMap[topic] = now.Unix() + int64(mqttClient.mqttRate)
		} else {
			if mqttClient.verbose {
				log.Printf("MQTT: skipped %s, rate to high", topic)
			}
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
	token := mqttClient.Publish(topic, byte(mqttQos), true, message)
	message = fmt.Sprintf("connected")
	if verbose {
		log.Printf("MQTT: push %s, message: %s", topic, message)
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
