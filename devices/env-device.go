package devices

import (
	"log"
	"time"

	auth "github.com/gcp-iotcore/simulated-devices/mqtt-auth"
	handler "github.com/gcp-iotcore/simulated-devices/mqtt-handler"
)

func StartEnvDevice(CertificatePath string, projectId string, privateKeyPath string,
	region string, registryId string, deviceId string, bridgeHost string, bridgePort string) error {
	log.Println("starting environment control device")
	log.Println("fetching TLS config")
	tlsConfig, err := auth.CreateTLSConfig(CertificatePath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("fetching jwt token")
	jwtToken, err := auth.JWTHandler(tlsConfig, projectId, privateKeyPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println(jwtToken)
	log.Println("configuring MQTT client")

	MQTTClient, err := handler.ConfigureClient(projectId, region, registryId, deviceId, bridgeHost, bridgePort, tlsConfig, jwtToken)

	for i := 0; i < 20; i++ {
		handler.PublishMessage(deviceId, MQTTClient, `{"message": "test"}`)
		time.Sleep(30 * time.Second)
	}

	return nil
}
