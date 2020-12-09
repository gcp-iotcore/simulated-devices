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
	log.Println(projectID)
	log.Println(region)
	log.Println(registryID)
	log.Println(deviceID)
	log.Println(bridgeHost)
	log.Println(bridgePort)
	log.Println(config)
	log.Println(jwtToken)
	opts := MQTT.NewClientOptions()
	broker := fmt.Sprintf("ssl://%v:%v", bridgeHost, bridgePort)
	log.Printf("Broker '%v'", broker)
	opts.AddBroker(broker)
	opts.SetClientID(clientID).SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12})
	opts.SetUsername("unused-device")
	opts.SetPassword(jwtToken)
	opts.SetProtocolVersion(4)
	// Incoming
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("[handler] Topic: %v\n", msg.Topic())
		fmt.Printf("[handler] Payload: %v\n", msg.Payload())
	})
	log.Println("[main] MQTT Client Connecting")
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token)
		log.Fatal(token.Error().Error())
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
