package devices

import (
	"bytes"
	"log"
	"time"

	httpHandler "github.com/gcp-iotcore/simulated-devices/http-handler"
	auth "github.com/gcp-iotcore/simulated-devices/mqtt-auth"
	handler "github.com/gcp-iotcore/simulated-devices/mqtt-handler"
)

func StartMQTTEnvDevice(CertificatePath string, projectId string, privateKeyPath string,
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

func StartHTTPEnvDevice(CertificatePath string, projectId string, privateKeyPath string,
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

	apiUrl := "projects/" + projectId + "/locations/" + region + "/registries/" + registryId + "/devices/" + deviceId + ":publishEvent"
	log.Println(apiUrl)
	for i := 0; i < 20; i++ {
		reqBody := bytes.NewBuffer([]byte(`{"binary_data": "eyJ1c2VyIjogImplcnJ5In0="}`))
		response, err := httpHandler.MakeHTTPCall("POST", apiUrl, jwtToken, reqBody)
		log.Println(response)
		log.Println(err)
		time.Sleep(30 * time.Second)
	}

	return nil
}
