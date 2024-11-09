package PMS7003

import (
	"fmt"
	"io"
	"log"

	"go.bug.st/serial"
)

type Mode string

const (
	ActiveMode  Mode = "ACTIVE"
	PassiveMode Mode = "PASSIVE"
)

type PMS7003Device struct {
	serialDevice     string
	mode             Mode
	serialConnection serial.Port
}

const startByte1 byte = 0x42
const startByte2 byte = 0x4D

// TODO: Make mode emum

func Open(serialDevice string, mode Mode) (device PMS7003Device, err error) {

	device = PMS7003Device{serialDevice, mode, nil}
	serialMode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	stream, err := serial.Open(serialDevice, serialMode)
	if err != nil {
		log.Println("unable to open serial connection.Please check device name and permissions are correct")
		return device, err
	}
	device.serialConnection = stream

	if mode == PassiveMode {
		setDevicePassive(&device, err)
		device.WakeUp()
	} else {
		setDeviceActive(&device, err)
		// wake up device in case it was left in sleep mode
		device.WakeUp()
	}

	return device, err
}

func (r *PMS7003Device) Read() (sensorValue PMS7003SensorValue, err error) {

	fmt.Printf("Device Initilized in %s Mode \n", r.mode)
	var dat []byte

	dat, err = readDataFromDevice(r)
	if err != nil {
		// could not read value so return the error
		return sensorValue, err
	}
	return parseSensorValueFromBytes(dat)

}

func readDataFromDevice(r *PMS7003Device) (dat []byte, err error) {
	dat = make([]byte, 32)
	_, err = io.ReadAtLeast(r.serialConnection, dat, 32)
	return dat, err
}

func (r *PMS7003Device) Close() {
	r.serialConnection.Close()
}

func (d *PMS7003Device) WakeUp() (err error) {

	var wakeUpCommandBytes []byte = []byte{startByte1, startByte2, 0xE4, 0x00, 0x01, 0x01, 0x74}

	x, x1 := d.serialConnection.Write(wakeUpCommandBytes)

	if x1 != nil && x != 7 {
		// TODO: split to error to differenciate device error vs not being able to write the full command
		err = fmt.Errorf("failed to wake up device")
	}

	// // wait for 500 mins as during testing it was found that device takes sometime for the fan to start
	// time.Sleep(5000 * time.Millisecond)
	return err
}

func (d *PMS7003Device) Sleep() (err error) {
	var sleepCommandBytes []byte = []byte{startByte1, startByte2, 0xE4, 0x00, 0x00, 0x01, 0x73}
	x, x1 := d.serialConnection.Write(sleepCommandBytes)
	if x1 != nil || x != 7 {
		// TODO: split to error to differenciate device error vs not being able to write the full command
		err = fmt.Errorf("failed to send sleep device")
	}
	return err
}

func setDeviceActive(device *PMS7003Device, err error) error {
	var activeModeCommand []byte = []byte{startByte1, startByte2, 0xE1, 0x00, 0x01, 0x01, 0x72}
	x, x1 := device.serialConnection.Write(activeModeCommand)
	if x1 != nil || x != 7 {
		// TODO: split to error to differenciate device error vs not being able to write the full command
		err = fmt.Errorf("failed to set device in active mode")
	}
	return err
}

func setDevicePassive(device *PMS7003Device, err error) error {
	var activeModeCommand []byte = []byte{startByte1, startByte2, 0xE1, 0x00, 0x01, 0x00, 0x71}
	x, x1 := device.serialConnection.Write(activeModeCommand)
	if x1 != nil || x != 7 {
		// TODO: split to error to differenciate device error vs not being able to write the full command
		err = fmt.Errorf("failed to set device in passive mode")
	}
	return err
}
