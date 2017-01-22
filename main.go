package main

import (
	"fmt"
	"github.com/dimorinny/adbaster/adb"
	"github.com/labstack/gommon/log"
)

func main() {
	config := adb.NewConfig(
		"192.168.99.100",
		32784,
		"adb",
	)
	client := adb.NewClient(config)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	info, err := client.DeviceInfo()
	if err != nil {
		log.Print(err)
	}

	fmt.Println(info)
}
