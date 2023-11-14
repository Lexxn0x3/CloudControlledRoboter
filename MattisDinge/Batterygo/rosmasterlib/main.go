package rosmasterlib

import (
	"log"

	"go.bug.st/serial"
)

type Rosmaster struct {
	port           serial.Port
	batteryVoltage float64
}

func NewRosmaster(comPort string, baudRate int) *Rosmaster {
	mode := &serial.Mode{
		BaudRate: baudRate,
	}
	port, err := serial.Open(comPort, mode)
	if err != nil {
		log.Fatal(err)
	}
	rm := &Rosmaster{
		port:           port,
		batteryVoltage: 0,
	}
	go rm.readSerial()
	return rm
}

func (rm *Rosmaster) Close() {
	rm.port.Close()
}

func (rm *Rosmaster) readSerial() {
	const head = 0xFF
	const deviceID = 0xFC

	for {
		buf := make([]byte, 1)
		_, err := rm.port.Read(buf)
		if err != nil {
			log.Println("Error reading from serial port:", err)
			continue
		}

		if buf[0] == head {
			_, err = rm.port.Read(buf) // Lies das nächste Byte (deviceID)
			if err != nil {
				continue
			}

			if buf[0] == deviceID-1 {
				header := make([]byte, 2)
				_, err := rm.port.Read(header) // Lies das Längen- und Typbyte
				if err != nil {
					continue
				}

				extLen := header[0]
				extType := header[1]
				checkSum := int(extLen) + int(extType)

				extData := make([]byte, extLen-2)
				var rxCheckNum int

				for i := range extData {
					_, err = rm.port.Read(buf)
					if err != nil {
						continue
					}

					if i == len(extData)-1 {
						rxCheckNum = int(buf[0])
					} else {
						checkSum += int(buf[0])
						extData[i] = buf[0]
					}
				}

				if byte(checkSum%256) == byte(rxCheckNum) {
					rm.parseData(extType, extData)
				} else {
					log.Println("Checksum error")
				}
			}
		}
	}
}

func (rm *Rosmaster) parseData(extType byte, data []byte) {
	// fmt.Println("extType: ", extType)
	// fmt.Println("data: ", data)
	if extType == 0x0A {
		rm.batteryVoltage = float64(data[6]) / 10
	}
}

func (rm *Rosmaster) GetBatteryVoltage() float64 {
	return rm.batteryVoltage
}
