package main

import (
	"fmt"

	devices "github.com/gcp-iotcore/simulated-devices/devices"
)

func main() {
	fmt.Println("main function")
	devices.StartHTTPEnvDevice()
}
