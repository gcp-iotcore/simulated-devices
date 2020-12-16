package devices

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math"
	"strings"
	"time"

	auth "github.com/gcp-iotcore/simulated-devices/mqtt-auth"
)

type AquaponicsDevice struct {
	DeviceType            string  `json:"device-type"`
	AquacultureWaterLevel float64 `json:"aquaculture-water-level"`
	CirculationPumpStatus string  `json:"circulation-pump-status"`
	ReservoirPumpStatus   string  `json:"reservoir-pump-status"`
	ReservoirWaterLevel   float64 `json:"reservoir-water-level"`
	TDS                   float64 `json:"tds"`
	WaterPh               float64 `json:"water-ph"`
	WaterTemperature      float64 `json:"water-temperature"`
}

func StartAquaponicsDevice() {
	log.Println("starting aquaponics control device")
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

	apiUrl := "projects/poc-cloudfaringpirates-797995/locations/us-central1/registries/aquaponics-poc/devices/aquaponics-device:publishEvent"
	log.Println(apiUrl)

	aquaponicsData := &AquaponicsDevice{}
	aquaponicsData.AquacultureWaterLevel = 80
	aquaponicsData.CirculationPumpStatus = "on"
	aquaponicsData.DeviceType = "aquaponics-device"
	aquaponicsData.ReservoirPumpStatus = "on"
	aquaponicsData.ReservoirWaterLevel = 40
	aquaponicsData.TDS = 350
	aquaponicsData.WaterPh = 7.6
	aquaponicsData.WaterTemperature = 27.5
	counter := 0

	for {

		aquaponicsData = handleAquaponicsMasterDevice(aquaponicsData)
		responseData := createAquaponicsJSON(aquaponicsData)
		response, err := deviceHTTPHandler(apiUrl, jwtToken, responseData)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println(response)

		time.Sleep(5 * time.Minute)
		counter++

		if counter == 200 {
			counter = 0
			jwtToken, err = auth.JWTHandler(tlsConfig, "poc-cloudfaringpirates-797995", "/simulated-devices/certs/rsa_private.pem")
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

}

func handleAquaponicsMasterDevice(data *AquaponicsDevice) *AquaponicsDevice {
	if data.CirculationPumpStatus == "on" {
		data.AquacultureWaterLevel = math.Round((data.AquacultureWaterLevel-2)*100) / 100
		data.ReservoirWaterLevel = math.Round((data.ReservoirWaterLevel+0.5)*100) / 100
		data.TDS = math.Round((data.TDS-5.7)*100) / 100
		data.WaterPh = math.Round((data.WaterPh-0.07)*100) / 100
		data.WaterTemperature = math.Round((data.WaterTemperature+0.01)*100) / 100
		if data.AquacultureWaterLevel <= 40 {
			data.ReservoirPumpStatus = "on"
			data.AquacultureWaterLevel = math.Round((data.AquacultureWaterLevel+5)*100) / 100
		}
		if data.AquacultureWaterLevel >= 80 {
			data.ReservoirPumpStatus = "off"
		}

		if data.TDS <= 250 {
			data.CirculationPumpStatus = "off"
			data.TDS = math.Round((data.TDS+1.2)*100) / 100
		}

		if data.WaterPh > 7.8 || data.WaterPh < 6.8 {
			data.WaterPh = 7.3
		}

	}
	if data.ReservoirPumpStatus == "on" {
		data.ReservoirWaterLevel = math.Round((data.ReservoirWaterLevel-0.8)*100) / 100
		data.AquacultureWaterLevel = math.Round((data.AquacultureWaterLevel+1.8)*100) / 100
		if data.AquacultureWaterLevel >= 80 {
			data.ReservoirPumpStatus = "off"
		}
	}
	return data
}

func createAquaponicsJSON(data *AquaponicsDevice) string {
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
