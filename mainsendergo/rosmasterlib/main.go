package rosmasterlib

import (
	"log"
	"time"

	"go.bug.st/serial"
)

type Rosmaster struct {
	port               serial.Port
	head               byte
	deviceID           byte
	complement         int
	batteryVoltage     float64
	ax, ay, az         float64
	gx, gy, gz         float64
	mx, my, mz         float64
	FUNC_MOTOR         byte
	FUNC_MOTION        byte
	FUNC_BEEP          byte
	FUNC_RGB           byte
	FUNC_RGB_EFFECT    byte
	__delay_time       float64
	__debug            bool
	BlockedHealthcheck bool
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
		port:            port,
		head:            0xFF,
		deviceID:        0xFC,
		complement:      0,
		batteryVoltage:  0,
		FUNC_MOTOR:      0x10,
		FUNC_MOTION:     0x12,
		FUNC_RGB:        0x05,
		FUNC_RGB_EFFECT: 0x06,
		FUNC_BEEP:       0x02,
		__delay_time:    0.002,
		__debug:         true,
	}

	//calc complement
	rm.complement = 257 - int(rm.deviceID)

	go rm.readSerial()
	return rm
}

func (rm *Rosmaster) Close() {
	rm.port.Close()
}

func (rm *Rosmaster) readSerial() {
	for {
		buf := make([]byte, 1)
		_, err := rm.port.Read(buf)
		if err != nil {
			log.Println("Warn reading from serial port: Maybe main program exited", err)
			return
		}

		if buf[0] == rm.head {
			_, err = rm.port.Read(buf) // Lies das nächste Byte (deviceID)
			if err != nil {
				continue
			}

			if buf[0] == rm.deviceID-1 {
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

func (rm *Rosmaster) writeSerial(data []byte) error {
	_, err := rm.port.Write(data)
	if err != nil {
		log.Println("Error writing to serial port:", err)
		return err
	}
	return nil
}

func (rm *Rosmaster) parseData(extType byte, data []byte) {
	// fmt.Println("extType: ", extType)
	// fmt.Println("data: ", data)
	if extType == 0x0A {
		rm.batteryVoltage = float64(data[6]) / 10
	}
	if extType == 0x0B {
		// Accelerometer
		accelRatio := 1 / 1000.0
		rm.ax = float64(int16(data[6])|int16(data[7])<<8) * accelRatio
		rm.ay = float64(int16(data[8])|int16(data[9])<<8) * accelRatio
		rm.az = float64(int16(data[10])|int16(data[11])<<8) * accelRatio

		// Gyroscope
		gyroRatio := 1 / 3754.9
		rm.gx = float64(int16(data[0])|int16(data[1])<<8) * gyroRatio
		rm.gy = float64(int16(data[2])|int16(data[3])<<8) * -gyroRatio
		rm.gz = float64(int16(data[4])|int16(data[5])<<8) * -gyroRatio

		// Magnitude
		magRatio := 1.0
		rm.mx = float64(int16(data[12])|int16(data[13])<<8) * magRatio
		rm.my = float64(int16(data[14])|int16(data[15])<<8) * magRatio
		rm.mz = float64(int16(data[16])|int16(data[17])<<8) * magRatio
	}

}

func (rm *Rosmaster) sum(data []byte, complement int) byte {
	sum := 0
	for _, b := range data {
		sum += int(b)
	}
	sum += complement
	return byte(sum & 0xff)
}

func (rm *Rosmaster) limitMotorValue(value int8) int8 {
	// Ensure the motor value is within the valid range of -100 to 100
	if value < -100 {
		return -100
	} else if value > 100 {
		return 100
	}
	return value
}

func (rm *Rosmaster) SetMotor(speed1, speed2, speed3, speed4 int8) {
	if rm.BlockedHealthcheck {
		return
	}
	tryCmd := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("---set_motor error!---")
			}
		}()

		tSpeedA := byte(rm.limitMotorValue(speed1))
		tSpeedB := byte(rm.limitMotorValue(speed2))
		tSpeedC := byte(rm.limitMotorValue(speed3))
		tSpeedD := byte(rm.limitMotorValue(speed4))

		cmd := []byte{
			rm.head,
			rm.deviceID,
			0x00,
			rm.FUNC_MOTOR,
			tSpeedA, tSpeedB, tSpeedC, tSpeedD,
		}

		cmd[2] = byte(len(cmd) - 1)
		checksum := byte(rm.sum(cmd, rm.complement) & 0xff)
		cmd = append(cmd, checksum)

		rm.writeSerial(cmd)

		if rm.__debug {
			log.Println("motor:", cmd)
		}

		time.Sleep(time.Duration(rm.__delay_time) * time.Second)
	}

	tryCmd()
}

func (rm *Rosmaster) SetCarMotion(vX, vY, vZ float64) {
	tryCmd := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("---set_car_motion error!---")
			}
		}()

		// Convert floats to int16 and then to byte array
		vxParms := int16(vX * 1000)
		vyParms := int16(vY * 1000)
		vzParms := int16(vZ * 1000)

		cmd := []byte{
			rm.head,
			rm.deviceID,
			0x00,
			rm.FUNC_MOTION,
			byte(vxParms & 0xFF),
			byte((vxParms >> 8) & 0xFF),
			byte(vyParms & 0xFF),
			byte((vyParms >> 8) & 0xFF),
			byte(vzParms & 0xFF),
			byte((vzParms >> 8) & 0xFF),
		}

		cmd[2] = byte(len(cmd) - 1)
		checksum := rm.sum(cmd, rm.complement)
		cmd = append(cmd, checksum)

		rm.writeSerial(cmd)

		if rm.__debug {
			log.Printf("motion: %v\n", cmd)
		}

		time.Sleep(time.Duration(rm.__delay_time) * time.Second)
	}

	tryCmd()
}

func (rm *Rosmaster) SetBeep(onTime int) {
	tryCmd := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("---set_beep error!---")
			}
		}()

		if onTime < 0 {
			log.Println("beep input error!")
			return
		}

		value := byte(onTime & 0xFF)            // Low byte
		valueHigh := byte((onTime >> 8) & 0xFF) // High byte

		cmd := []byte{
			rm.head,
			rm.deviceID,
			0x05,
			rm.FUNC_BEEP,
			value,
			valueHigh,
		}

		cmd[2] = byte(len(cmd) - 1)
		checksum := byte(rm.sum(cmd, rm.complement) & 0xff)
		cmd = append(cmd, checksum)

		rm.writeSerial(cmd)

		if rm.__debug {
			log.Println("beep:", cmd)
		}

		time.Sleep(time.Duration(rm.__delay_time) * time.Second)
	}

	tryCmd()
}

func (rm *Rosmaster) SetColorfulLamps(ledID, red, green, blue byte) {
	tryCmd := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("---set_colorful_lamps error!---")
			}
		}()

		cmd := []byte{
			rm.head,
			rm.deviceID,
			0x00,
			rm.FUNC_RGB,
			ledID,
			red,
			green,
			blue,
		}

		cmd[2] = byte(len(cmd) - 1)
		checksum := rm.sum(cmd, rm.complement)
		cmd = append(cmd, checksum)

		err := rm.writeSerial(cmd)
		if err != nil {
			log.Println("Error writing rgb command to serial port:", err)
			return
		}

		if rm.__debug {
			log.Printf("rgb: %v\n", cmd)
		}

		time.Sleep(time.Duration(rm.__delay_time) * time.Second)
	}

	tryCmd()
}

func (rm *Rosmaster) SetColorfulEffect(effect, speed, parm byte) {
	tryCmd := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("---set_colorful_effect error!---")
			}
		}()

		cmd := []byte{
			rm.head,
			rm.deviceID,
			0x00,
			rm.FUNC_RGB_EFFECT,
			effect,
			speed,
			parm,
		}

		cmd[2] = byte(len(cmd) - 1)
		checksum := rm.sum(cmd, rm.complement)
		cmd = append(cmd, checksum)

		err := rm.writeSerial(cmd)
		if err != nil {
			log.Println("Error writing rgb effect command to serial port:", err)
			return
		}

		if rm.__debug {
			log.Printf("rgb_effect: %v\n", cmd)
		}

		time.Sleep(time.Duration(rm.__delay_time) * time.Second)
	}

	tryCmd()
}

func (rm *Rosmaster) GetBatteryVoltage() float64 {
	return rm.batteryVoltage
}

func (rm *Rosmaster) GetGyroscope() (float64, float64, float64) {
	return rm.gx, rm.gy, rm.gz
}

func (rm *Rosmaster) GetAcceleration() (float64, float64, float64) {
	return rm.ax, rm.ay, rm.az
}

func (rm *Rosmaster) GetMagnitude() (float64, float64, float64) {
	return rm.mx, rm.my, rm.mz
}
