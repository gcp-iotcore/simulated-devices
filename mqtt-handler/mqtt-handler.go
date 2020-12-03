package mqtthandler

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func ConfigureClient(projectID string, region string, registryID string, deviceID string,
	bridgeHost string, bridgePort string, config *tls.Config, jwtToken string) (MQTT.Client, error) {
	clientID := fmt.Sprintf("projects/%v/locations/%v/registries/%v/devices/%v",
		projectID,
		region,
		registryID,
		deviceID,
	)
	log.Println("Creating MQTT Client Options")
	opts := MQTT.NewClientOptions()
	broker := fmt.Sprintf("ssl://%v:%v", bridgeHost, bridgePort)
	log.Printf("Broker '%v'", broker)
	opts.AddBroker(broker)
	opts.SetClientID(clientID).SetTLSConfig(config)
	opts.SetUsername("unused")
	opts.SetPassword(jwtToken)
	// Incoming
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("[handler] Topic: %v\n", msg.Topic())
		fmt.Printf("[handler] Payload: %v\n", msg.Payload())
	})
	log.Println("[main] MQTT Client Connecting")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
		return nil, token.Error()
	}
	return client, nil
}

func PublishMessage(deviceID string, client MQTT.Client, message string) {
	topic := struct {
		config    string
		telemetry string
	}{
		config:    fmt.Sprintf("/devices/%v/config", deviceID),
		telemetry: fmt.Sprintf("/devices/%v/events", deviceID),
	}

	log.Println("Creating Subscription")
	client.Subscribe(topic.config, 0, nil)
	log.Println("Publishing Message", message)
	token := client.Publish(
		topic.telemetry,
		0,
		false,
		fmt.Sprintf(message))
	token.WaitTimeout(5 * time.Second)
}
