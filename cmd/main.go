package main

import (
	"fmt"
	"time"

	PMS7003 "github.com/shivasaxena/PMS7003"
)

func main() {

	device, err := PMS7003.Open("/dev/ttyAMA0", PMS7003.PassiveMode)

	if err != nil {
		panic(err)
	}

	defer device.Close()

	device.WakeUp()
	time.Sleep(10 * time.Second)
	for i := 0; i < 5; i++ {

		value, err := device.Read()

		if err != nil {
			panic(err)
		}

		fmt.Println(value.PM25Atmospheric)
	}

	device.Sleep()

}
