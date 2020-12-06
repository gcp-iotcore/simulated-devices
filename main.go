package main

import (
	"flag"
	"fmt"

	devices "github.com/gcp-iotcore/simulated-devices/devices"
)

var (
	deviceID = flag.String("device", "", "Cloud IoT Core Device ID")
	bridge   = struct {
		host *string
		port *string
	}{
		flag.String("mqtt_host", "mqtt.googleapis.com", "MQTT Bridge Host"),
		flag.String("mqtt_port", "8883", "MQTT Bridge Port"),
	}
	projectID  = flag.String("project", "", "GCP Project ID")
	registryID = flag.String("registry", "", "Cloud IoT Registry ID (short form)")
	region     = flag.String("region", "", "GCP Region")
	certsCA    = flag.String("ca_certs", "", "Download https://pki.google.com/roots.pem")
	privateKey = flag.String("private_key", "", "Path to private key file")
)

func main() {
	fmt.Println("main function")
	flag.Parse()
	devices.StartEnvDevice(*certsCA, *projectID, *privateKey, *region, *registryID, *deviceID, *bridge.host, *bridge.port)
}
