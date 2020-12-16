package devices

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"

	auth "github.com/gcp-iotcore/simulated-devices/mqtt-auth"
)

type MushroomsDeviceData struct {
	DeviceType      string  `json:"device-type"`
	LightingLevel   float64 `json:"lighting-level"`
	SprinklerStatus string  `json:"sprinkler-status"`
}

type EarthwormsDeviceData struct {
	DeviceType      string  `json:"device-type"`
	LightingLevel   float64 `json:"lighting-level"`
	SoilPh          float64 `json:"soil-ph"`
	SoilTemperature float64 `json:"soil-temperature"`
}

func StartMushroomsDevice() {
	log.Println("starting mushrooms control device")
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

	apiUrl := "projects/poc-cloudfaringpirates-797995/locations/us-central1/registries/aquaponics-poc/devices/mushrooms-device:publishEvent"
	log.Println(apiUrl)

	mushroomsData := &MushroomsDeviceData{}
	mushroomsData.DeviceType = "mushrooms-device"
	mushroomsData.LightingLevel = 6.5
	mushroomsData.SprinklerStatus = "on"

	counter1 := 0
	counter2 := 0

	for {

		mushroomsData = handleMushroomsMasterDevice(mushroomsData, counter1, counter2)
		responseData := createMushroomsJSON(mushroomsData)
		response, err := deviceHTTPHandler(apiUrl, jwtToken, responseData)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(response)
		time.Sleep(5 * time.Minute)
		counter1++
		counter2++
		if counter1 > 9 {
			counter1 = 0
		}

		if counter2 > 5 {
			counter2 = 0
		}

	}
}

func StartEarthwormsDevice() {
	log.Println("starting earthworms control device")
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

	apiUrl := "projects/poc-cloudfaringpirates-797995/locations/us-central1/registries/aquaponics-poc/devices/earthworms-device:publishEvent"
	log.Println(apiUrl)
	earthWormsData := &EarthwormsDeviceData{}
	earthWormsData.DeviceType = "earthworms-device"
	earthWormsData.LightingLevel = 6.5
	earthWormsData.SoilPh = 7.7
	earthWormsData.SoilTemperature = 27.8

	counter := 0

	for {
		earthWormsData = handleEarthwormsMasterDevice(earthWormsData, counter)
		responseData := createEarthWormsJSON(earthWormsData)
		response, err := deviceHTTPHandler(apiUrl, jwtToken, responseData)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(response)
		time.Sleep(5 * time.Minute)
		counter++
		if counter > 9 {
			counter = 0
		}
	}

}

func handleMushroomsMasterDevice(data *MushroomsDeviceData, counter1 int, counter2 int) *MushroomsDeviceData {

	var lux = [10]float64{5.5, 5.8, 5.9, 6.3, 6.5, 6.8, 6.9, 7.0, 7.1, 7.4}
	data.LightingLevel = lux[counter1]

	if counter2 == 5 {
		if data.SprinklerStatus == "on" {
			data.SprinklerStatus = "off"
		} else {
			data.SprinklerStatus = "on"
		}
	}

	return data
}

func handleEarthwormsMasterDevice(data *EarthwormsDeviceData, counter int) *EarthwormsDeviceData {

	var lux = [10]float64{5.5, 5.8, 5.9, 6.3, 6.5, 6.8, 6.9, 7.0, 7.1, 7.4}
	var temp = [10]float64{26, 26.7, 26.9, 27.1, 27.5, 27.8, 27.9, 28.4, 28.6, 29.5}
	var ph = [10]float64{6.9, 7.05, 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.4, 7.2}
	data.LightingLevel = lux[counter]
	data.SoilTemperature = temp[counter]
	data.SoilPh = ph[counter]

	return data
}

func createMushroomsJSON(data *MushroomsDeviceData) string {
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

func createEarthWormsJSON(data *EarthwormsDeviceData) string {
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
