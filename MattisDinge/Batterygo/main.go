package main

import (
	"fmt"
	"time"

	"github.com/mattzi/gobattery/rosmasterlib"
	// Ersetzen Sie 'your_project_name' durch den Namen Ihres Projekts
)

func main() {
	// Ersetzen Sie "/dev/ttyUSB0" durch Ihren seriellen Port und 115200 durch Ihre Baudrate
	rm := rosmasterlib.NewRosmaster("/dev/myserial", 115200)
	defer rm.Close()

	for {
		voltage := rm.GetBatteryVoltage()
		fmt.Printf("Aktuelle Batteriespannung: %v V\n", voltage)
		time.Sleep(2 * time.Second)
	}
}
