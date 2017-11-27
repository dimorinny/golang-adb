package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dimorinny/adbaster"
	"github.com/dimorinny/adbaster/adb"
	"github.com/dimorinny/adbaster/model"
	"log"
	"os"
)

func main() {
	testPackage := os.Args[1]
	applicationPackage := os.Args[2]
	runner := os.Args[3]
	testClass := os.Args[4]

	client := createClient()
	device := getFirstConnectedDevice(client)

	installApplications(client, device, "avito.apk", "avito-test.apk")

	runTests(
		client,
		device,
		testPackage,
		applicationPackage,
		runner,
		testClass,
	)
}

func createClient() adbaster.Client {
	config := adb.NewConfig(
		"adb",
		"\n",
	)
	return adb.NewClient(config, true)
}

func getFirstConnectedDevice(client adbaster.Client) model.DeviceIdentifier {
	identifiers, err := client.Devices()
	if err != nil {
		log.Fatal(err)
	}

	return identifiers[0]
}

func installApplications(client adbaster.Client, device model.DeviceIdentifier, applications ...string) {
	for _, application := range applications {
		err := client.Install(device, application)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func listenLogcat(client adbaster.Client, device model.DeviceIdentifier) {
	logcatStream, err := client.Logcat(device)
	if err != nil {
		log.Fatal(err)
	}

	for line := range logcatStream {
		fmt.Println(line)
	}
}

func runTests(
	client adbaster.Client,
	device model.DeviceIdentifier,
	testPackage,
	applicationPackage,
	runner,
	testClass string,

) {
	eventStream, instrumentationOutput, err := client.RunInstrumentationTests(
		device,
		model.InstrumentationParams{
			TestPackage: testPackage,
			Runner:      runner,
			Arguments: model.InstrumentationArguments{
				"testType":                            "FIREBASE",
				"fileStorageHost":                     "erc20.xyz",
				"fileStorageAccessKey":                "access",
				"fileStorageSecretKey":                "secret",
				"componentTestFlakyFilterIterations":  "2",
				"componentTestTakeScreenshots":        "false",
				"allureReportForInstrumentationTests": "false",
				"componentTestFlakyDebug":             "false",
				"class":                               testClass,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// getting original instrumentation output
	for line := range instrumentationOutput {
		fmt.Println(line)
	}

	// getting test results
	for event := range eventStream {
		spew.Dump(event)
	}

	err = client.ClearApplicationData(device, applicationPackage)
	if err != nil {
		log.Fatal(err)
	}
}
