package rplidar

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"go.bug.st/serial"
)

const (
	SyncByte      = 0xA5
	SyncByte2     = 0x5A
	GetInfoByte   = 0x50
	StartScanByte = 0x20
	StopByte      = 0x25
	ResetByte     = 0x40
	// ... other constants as needed
)

type RPLidar struct {
	serialPort serial.Port
	portName   string
	baudrate   int
	timeout    time.Duration
}

type Info struct {
	Model        byte
	Firmware     [2]byte
	Hardware     byte
	SerialNumber []byte
}

type Measurement struct {
	Quality  byte
	Angle    float32
	Distance float32
}

func NewRPLidar(portName string, baudrate int, timeout time.Duration) *RPLidar {
	return &RPLidar{
		portName: portName,
		baudrate: baudrate,
		timeout:  timeout,
	}
}

func (rpl *RPLidar) parseRawScanData(data []byte) (bool, int, float32, float32) {
	s := data[0]&0x01 > 0
	abS := (data[0]>>1)&0x01 > 0
	quality := int(data[0]) >> 2
	if s == abS {
		log.Println("Data check bit failed while parsing raw scan data")
	}
	c := data[1] & 0x01
	if c != 1 {
		log.Println("Check bit was not 1")
	}
	angle := float32(int(data[1])>>1+int(data[2])<<7) / 64.0
	distance := float32(int(data[3])+(int(data[4])<<8)) / 4.0
	return s, quality, angle, distance
}

func (lidar *RPLidar) Connect() error {
	mode := &serial.Mode{
		BaudRate: lidar.baudrate,
		// ... other serial port settings
	}
	port, err := serial.Open(lidar.portName, mode)
	if err != nil {
		return err
	}
	lidar.serialPort = port
	return nil
}

func (lidar *RPLidar) Disconnect() error {
	if lidar.serialPort != nil {
		return lidar.serialPort.Close()
	}
	return nil
}

func (lidar *RPLidar) sendCommand(cmd byte) error {
	_, err := lidar.serialPort.Write([]byte{SyncByte, cmd})
	return err
}

func (lidar *RPLidar) readResponse(size int) ([]byte, error) {
	deadline := time.Now().Add(lidar.timeout)
	buffer := make([]byte, size)
	totalRead := 0

	for totalRead < size {
		if time.Now().After(deadline) {
			return nil, errors.New("read response: timeout reached")
		}

		n, err := lidar.serialPort.Read(buffer[totalRead:])
		if err != nil {
			if err == io.EOF {
				continue
			}
			return nil, err
		}

		totalRead += n
	}

	return buffer, nil
}

func (lidar *RPLidar) GetInfo() (Info, error) {
	err := lidar.sendCommand(GetInfoByte)
	if err != nil {
		return Info{}, err
	}

	buf, err := lidar.readResponse(20)
	if err != nil {
		return Info{}, err
	}

	return Info{
		Model:        buf[0],
		Firmware:     [2]byte{buf[2], buf[1]},
		Hardware:     buf[3],
		SerialNumber: buf[4:],
	}, nil
}

func (lidar *RPLidar) StartScan() (int, error) {
	if err := lidar.StartMotor(); err != nil {
		return 0, fmt.Errorf("start motor error: %w", err)
	}
	//lidar.readResponse(1)
	lidar.serialPort.ResetInputBuffer()

	err := lidar.sendCommand(StartScanByte)
	if err != nil {
		return 0, fmt.Errorf("send start scan command error: %w", err)
	}

	// Read descriptor to get asize
	descriptor, err := lidar.readResponse(7) // Descriptor is 7 bytes
	if err != nil {
		return 0, fmt.Errorf("read descriptor error: %w", err)
	}

	fmt.Println(descriptor)
	if len(descriptor) != 7 || descriptor[0] != SyncByte || descriptor[1] != SyncByte2 {
		return 0, errors.New("invalid descriptor received")
	}

	asize := int(descriptor[2]) // Extract asize from descriptor

	// Discard initial data based on asize
	_, err = lidar.readResponse(asize)
	if err != nil {
		return 0, fmt.Errorf("read initial data error: %w", err)
	}

	return asize, nil
}

func (lidar *RPLidar) StopScan() error {
	err := lidar.sendCommand(StopByte)
	if err != nil {
		return fmt.Errorf("send stop scan command error: %w", err)
	}

	if err := lidar.StopMotor(); err != nil {
		return fmt.Errorf("stop motor error: %w", err)
	}

	return nil
}

func (lidar *RPLidar) StartMotor() error {
	// Start the motor (set DTR to false)
	err := lidar.serialPort.SetDTR(false)
	if err != nil {
		return fmt.Errorf("start motor setDTR error: %w", err)
	}
	return nil
}

func (lidar *RPLidar) StopMotor() error {
	// Stop the motor (set DTR to true)
	err := lidar.serialPort.SetDTR(true)
	if err != nil {
		return fmt.Errorf("stop motor setDTR error: %w", err)
	}
	return nil
}

func (lidar *RPLidar) IterMeasurements() (<-chan Measurement, error) {
	asize, err := lidar.StartScan()
	if err != nil {
		return nil, fmt.Errorf("start scan error: %w", err)
	}

	ch := make(chan Measurement)

	go func() {
		defer close(ch)
		for {
			data, err := lidar.readResponse(asize) // Use asize here
			if err != nil {
				if err == io.EOF {
					continue
				}
				log.Printf("Error reading measurement: %v", err)
				return
			}

			if len(data) != asize {
				log.Printf("Unexpected data length: %d, expected: %d", len(data), asize)
				continue
			}

			newScan, quality, angle, distance := lidar.parseRawScanData(data)
			if newScan {
				// Handle new scan logic if needed
			}

			if quality > 0 && distance > 0 {
				ch <- Measurement{Quality: byte(quality), Angle: angle, Distance: distance}
			}
		}
	}()

	return ch, nil
}
