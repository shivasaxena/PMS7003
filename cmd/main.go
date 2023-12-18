package main

import (
	"fmt"

	PMS7003 "github.com/shivasaxena/PMS7003"
)

func main() {

	device, err := PMS7003.Open("/dev/ttyAMA0", PMS7003.ActiveMode)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {

		value, err := device.Read()

		if err != nil {
			panic(err)
		}

		fmt.Println(value.PM25Atmospheric)
	}

	defer device.Close()

}
