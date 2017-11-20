package adbaster

import (
	"github.com/dimorinny/adbaster/adb"
	"github.com/dimorinny/adbaster/model"
	"log"
	"testing"
	"fmt"
)

func TestMain1(t *testing.T) {
	client := createClient()
	device := getFirstConnectedDevice(client)

	installApplications(client, device, "avito.apk", "avito-test.apk")

	runTests(client, device)
}

func runTests(client Client, device model.DeviceIdentifier) {
	result, err := client.RunInstrumentationTests(
		device,
		model.InstrumentationParams{
			TestPackage: "com.avito.android.dev.test",
			Runner:      "com.avito.android.runner.AvitoInstrumentTestRunner",
			Arguments: model.InstrumentationArguments{
				"testType":                            "FIREBASE",
				"fileStorageHost":                     "erc20.xyz",
				"fileStorageAccessKey":                "access",
				"fileStorageSecretKey":                "secret",
				"componentTestFlakyFilterIterations":  "2",
				"componentTestTakeScreenshots":        "false",
				"allureReportForInstrumentationTests": "false",
				"componentTestFlakyDebug":             "false",
				"class":                               "com.avito.android.module.edit_profile.AvatarEditProfileTest",
			},
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	for message := range result {
		fmt.Println(message)
		fmt.Println("========")
	}
}

func createClient() Client {
	config := adb.NewConfig("adb")
	return adb.NewClient(config, true)
}

func getFirstConnectedDevice(client Client) model.DeviceIdentifier {
	identifiers, err := client.Devices()
	if err != nil {
		log.Fatal(err)
	}

	return identifiers[0]
}

func installApplications(client Client, device model.DeviceIdentifier, applications ...string) {
	for _, application := range applications {
		err := client.Install(device, application)
		if err != nil {
			log.Fatal(err)
		}
	}
}
