package main

import (
	"github.com/dimorinny/adbaster/adb"
	"github.com/labstack/gommon/log"
	"fmt"
)

func main() {
	config := adb.NewConfig(
		"192.168.99.100",
		32800,
		"adb",
	)
	client := adb.NewClient(config)

	fmt.Println("Connecting to device...")

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Getting device info...")

	info, err := client.DeviceInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Installing applications...")

	err = client.Install(
		"/Users/damerkurev/Documents/Projects/Android/Protools/avito-services-pro/services/build/outputs/apk/services-2.3.0-debug-androidTest.apk",
		"/data/local/tmp/com.avito.services.debug.test",
	)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Install(
		"/Users/damerkurev/Documents/Projects/Android/Protools/avito-services-pro/services/build/outputs/apk/services-2.3.0-debug.apk",
		"/data/local/tmp/com.avito.services.debug",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Run instrumentation tests...")

	err = client.RunInstrumentation(
		"com.avito.services.debug.test",
		"ru.avito.services.ServicesTestRunner",
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info)
}
