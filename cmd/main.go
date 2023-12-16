package main

import (
	"fmt"

	PMS7003 "github.com/shivasaxena/PMS7003/pkg"
)

func main() {

	device, err := PMS7003.Open("/dev/ttyAMA0", "Active")

	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println(device.Read())
	}

	defer device.Close()

}
