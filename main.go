package main

import (
	"fmt"
	"github.com/dimorinny/adbaster/adb"
	"github.com/labstack/gommon/log"
)

func main() {
	config := adb.NewConfig(
		"192.168.99.100",
		32814,
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

	result, err := client.RunInstrumentationTests(
		adb.InstrumentationParams{
			From:      "com.avito.services.debug.test",
			Runner:    "ru.avito.services.ServicesTestRunner",
			TestClass: "ru.avito.services.test.blacklist.BlacklistTest",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	fmt.Println(info)
}
