package mqtthandler

import (
	"crypto/tls"
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func ConfigureClient(projectID string, region string, registryID string, deviceID string,
	bridgeHost string, bridgePort string, config *tls.Config, jwtToken string) {
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
}
