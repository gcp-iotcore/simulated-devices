package devices

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	httpHandler "github.com/gcp-iotcore/simulated-devices/http-handler"
	auth "github.com/gcp-iotcore/simulated-devices/mqtt-auth"
	handler "github.com/gcp-iotcore/simulated-devices/mqtt-handler"
)

type EnvDevice struct {
	DeviceType        string  `json:"device-type"`
	EnvType           string  `json:"env-type"`
	FanStatusExhaust  string  `json:"fan-status-exhaust"`
	FanStatusInternal string  `json:"fan-status-internal"`
	RelativeHumidity  float64 `json:"relative-humidity"`
	RoomTemp          float64 `json:"room-temp"`
}

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

func StartHTTPEnvDevice() {
	log.Println("starting environment control device")
	log.Println("fetching TLS config")
	tlsConfig, err := auth.CreateTLSConfig("/simulated-devices/certs/roots.pem")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("fetching jwt token")
	jwtToken, err := auth.JWTHandler(tlsConfig, "poc-cloudfaringpirates-797995", "/simulated-devices/certs/rsa_private.pem")
	if err != nil {
		log.Fatal(err)
	}

	apiUrl := "projects/poc-cloudfaringpirates-797995/locations/us-central1/registries/aquaponics-poc/devices/env-device:publishEvent"
	log.Println(apiUrl)
	// for i := 0; i < 20; i++ {
	// 	reqBody := bytes.NewBuffer([]byte(`{"binary_data": "eyJ1c2VyIjogImplcnJ5In0="}`))
	// 	response, err := httpHandler.MakeHTTPCall("POST", apiUrl, jwtToken, reqBody)
	// 	log.Println(response)
	// 	log.Println(err)
	// 	time.Sleep(30 * time.Second)
	// }

	aquaponicsData := &EnvDevice{}
	aquaponicsData.DeviceType = "env-device"
	aquaponicsData.EnvType = "aquaponics"
	aquaponicsData.RoomTemp = 30.5
	aquaponicsData.RelativeHumidity = 90.4
	aquaponicsData.FanStatusExhaust = "on"
	aquaponicsData.FanStatusInternal = "on"

	mushroomsData := &EnvDevice{}
	mushroomsData.DeviceType = "env-device"
	mushroomsData.EnvType = "mushrooms"
	mushroomsData.RoomTemp = 29.5
	mushroomsData.RelativeHumidity = 93.4
	mushroomsData.FanStatusExhaust = "on"
	mushroomsData.FanStatusInternal = "on"

	earthwormsData := &EnvDevice{}
	earthwormsData.DeviceType = "env-device"
	earthwormsData.EnvType = "earthworms"
	earthwormsData.RoomTemp = 32.5
	earthwormsData.RelativeHumidity = 90.4
	earthwormsData.FanStatusExhaust = "on"
	earthwormsData.FanStatusInternal = "on"

	counter := 0

	for {

		aquaponicsData = handleAquaponicsDevice(aquaponicsData)
		responseData := createEnvJSON(aquaponicsData)
		response, err := deviceHTTPHandler(apiUrl, jwtToken, responseData)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(response)

		time.Sleep(30 * time.Second)

		mushroomsData = handleMushroomsDevice(mushroomsData)
		responseData = createEnvJSON(mushroomsData)
		response, err = deviceHTTPHandler(apiUrl, jwtToken, responseData)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(response)

		time.Sleep(30 * time.Second)

		earthwormsData = handleEarthwormsDevice(earthwormsData)
		responseData = createEnvJSON(earthwormsData)
		response, err = deviceHTTPHandler(apiUrl, jwtToken, responseData)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(response)
		time.Sleep(3 * time.Minute)
		counter++
		log.Println("counter ", counter)
		if counter == 300 {
			counter = 0
			jwtToken, err = auth.JWTHandler(tlsConfig, "poc-cloudfaringpirates-797995", "/simulated-devices/certs/rsa_private.pem")
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

}

func deviceHTTPHandler(apiUrl string, jwtToken string, responseData string) (*http.Response, error) {
	reqBody := bytes.NewBuffer([]byte(responseData))
	response, err := httpHandler.MakeHTTPCall("POST", apiUrl, jwtToken, reqBody)
	return response, err
}

func createEnvJSON(data *EnvDevice) string {
	apiData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(apiData))

	responseData := `{"binary_data": "replace"}`

	encodedData := base64.StdEncoding.EncodeToString(apiData)

	responseData = strings.Replace(responseData, "replace", encodedData, -1)

	return responseData
}

func handleAquaponicsDevice(data *EnvDevice) *EnvDevice {

	if data.FanStatusExhaust == "on" {
		data.RelativeHumidity = math.Round((data.RelativeHumidity-0.8)*100) / 100
		data.RoomTemp = math.Round((data.RoomTemp-0.2)*100) / 100
		if data.RelativeHumidity <= 85 {
			data.FanStatusExhaust = "off"
		}
	} else {
		data.RelativeHumidity = math.Round((data.RelativeHumidity+0.4)*100) / 100
		if data.RelativeHumidity >= 90 {
			data.FanStatusExhaust = "on"
		}
	}

	if data.FanStatusInternal == "on" {
		data.RoomTemp = math.Round((data.RoomTemp-0.5)*100) / 100
		if data.RoomTemp <= 27 {
			data.FanStatusInternal = "off"
		}
	} else {
		data.RoomTemp = math.Round((data.RoomTemp+0.1)*100) / 100
		if data.RoomTemp >= 31 {
			data.FanStatusInternal = "on"
		}
	}

	return data
}

func handleMushroomsDevice(data *EnvDevice) *EnvDevice {

	if data.FanStatusExhaust == "on" {
		data.RelativeHumidity = math.Round((data.RelativeHumidity-0.7)*100) / 100
		data.RoomTemp = math.Round(data.RoomTemp - 0.3)
		if data.RelativeHumidity <= 88 {
			data.FanStatusExhaust = "off"
		}
	} else {
		data.RelativeHumidity = math.Round((data.RelativeHumidity+0.8)*100) / 100
		if data.RelativeHumidity >= 94 {
			data.FanStatusExhaust = "on"
		}
	}

	if data.FanStatusInternal == "on" {
		data.RoomTemp = math.Round((data.RoomTemp-1.2)*100) / 100
		if data.RoomTemp <= 25 {
			data.FanStatusInternal = "off"
		}
	} else {
		data.RoomTemp = math.Round((data.RoomTemp+0.2)*100) / 100
		if data.RoomTemp >= 29 {
			data.FanStatusInternal = "on"
		}
	}

	return data
}

func handleEarthwormsDevice(data *EnvDevice) *EnvDevice {

	if data.FanStatusExhaust == "on" {
		data.RelativeHumidity = math.Round((data.RelativeHumidity-0.8)*100) / 100
		data.RoomTemp = math.Round((data.RoomTemp-0.1)*100) / 100
		if data.RelativeHumidity <= 88 {
			data.FanStatusExhaust = "off"
		}
	} else {
		data.RelativeHumidity = math.Round((data.RelativeHumidity+0.4)*100) / 100
		if data.RelativeHumidity >= 93 {
			data.FanStatusExhaust = "on"
		}
	}

	if data.FanStatusInternal == "on" {
		data.RoomTemp = math.Round((data.RoomTemp-0.9)*100) / 100
		if data.RoomTemp <= 29 {
			data.FanStatusInternal = "off"
		}
	} else {
		data.RoomTemp = math.Round((data.RoomTemp+0.2)*100) / 100
		if data.RoomTemp >= 32 {
			data.FanStatusInternal = "on"
		}
	}

	return data
}
