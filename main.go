package main

import (
	"fmt"

	auth "github.com/gcp-iotcore/simulated-devices/mqtt-auth"
)

func main() {
	fmt.Println("main function")
	auth.Dummy()
}
