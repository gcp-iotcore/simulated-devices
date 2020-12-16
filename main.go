package main

import (
	"fmt"
	"time"

	devices "github.com/gcp-iotcore/simulated-devices/devices"
)

func main() {
	fmt.Println("main function")
	fmt.Println("starting env device")
	go devices.StartHTTPEnvDevice()

	time.Sleep(20 * time.Second)
	fmt.Println("starting aquaponics device")
	go devices.StartAquaponicsDevice()

	time.Sleep(30 * time.Second)
	fmt.Println("starting earthworms device")
	go devices.StartEarthwormsDevice()
	time.Sleep(15 * time.Second)
	fmt.Println("starting mushrooms device")
	devices.StartMushroomsDevice()
}
