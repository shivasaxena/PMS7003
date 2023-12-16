package PMS7003

import (
	"bytes"
	"encoding/binary"
	"io"

	"go.bug.st/serial"
)

type Mode string

const (
	Active  Mode = "ACTIVE"
	Passive Mode = "PASSIVE"
)

type PMS7003Device struct {
	serialDevice     string
	mode             string // ACTIVE or PASSIVE
	serialConnection serial.Port
}

const startByte1 byte = 0x42
const startByte2 byte = 0x4D

var sleepCommandBytes []byte = []byte{startByte1, startByte2, 0xE4, 0x00, 0x00, 0x01, 0x73}
var wakeUpCommandBytes []byte = []byte{startByte1, startByte2, 0xE4, 0x00, 0x01, 0x01, 0x74}

// TODO: Make mode emum

func Open(serialDevice string, mode string) (device PMS7003Device, err error) {

	device = PMS7003Device{serialDevice, mode, nil}
	serialMode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	stream, err := serial.Open(serialDevice, serialMode)
	if err != nil {
		panic("Serial Port did not open")
	}
	device.serialConnection = stream

	return device, err
}

func (r *PMS7003Device) Read() (sensorValue PMS7003SensorValue, err error) {

	dat := make([]byte, 32)
	_, err = io.ReadAtLeast(r.serialConnection, dat, 32)
	if err != nil {
		panic(err)
	}

	return parseSensorValueFromBytes(dat)

}

func (r *PMS7003Device) Close() {
	r.serialConnection.Close()
}

func parseSensorValueFromBytes(rawBytes []byte) (sensorValue PMS7003SensorValue, err error) {
	// The following document lists the trasmission format under
	// Appendix I：PMS7003 transport protocol-Active Mode
	// http://download.kamami.pl/p564008-PMS7003%20series%20data%20manua_English_V2.5.pdf

	err = binary.Read(bytes.NewBuffer(rawBytes[2:]), binary.BigEndian, &sensorValue)

	checksum := uint16(0)
	for i := 0; i < 30; i++ {
		checksum += uint16(rawBytes[i])
	}

	return sensorValue, err
}
