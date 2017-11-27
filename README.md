## Description

Library for communication with Android device through Adb. The main purpose of this library is running Android instrumentation tests and parsing it's results.

## Usage

The main unit of this library is `adbaster.Client` interface, that looks like this:

```go
type Client interface {
	Devices() ([]model.DeviceIdentifier, error)
	DeviceInfo(device model.DeviceIdentifier) (*model.Device, error)
	Push(device model.DeviceIdentifier, from, to string) error
	Pull(device model.DeviceIdentifier, from, to string) error
	Install(device model.DeviceIdentifier, from string) error
	ClearApplicationData(device model.DeviceIdentifier, applicationPackage string) error
	Logcat(device model.DeviceIdentifier) (
		<-chan string,
		error,
	)
	RunInstrumentationTests(device model.DeviceIdentifier, params model.InstrumentationParams) (
		<-chan instrumentation.Event,
		<-chan string,
		error,
	)
}
```

This library has one implementation of this interface, that uses Adb binary. You can create it like this:

```go
func createClient() adbaster.Client {
	config := adb.NewConfig(
		"adb",
		"\n",
	)
	return adb.NewClient(config, true)
}
```

After that you can get connected device for further communication. For example, you can get first device like this:

```go
func getFirstConnectedDevice(client adbaster.Client) model.DeviceIdentifier {
	identifiers, err := client.Devices()
	if err != nil {
		log.Fatal(err)
	}

	return identifiers[0]
}
```

Now, you can install applications to your device:

```go
func installApplications(client adbaster.Client, device model.DeviceIdentifier, applications ...string) {
	for _, application := range applications {
		err := client.Install(device, application)
		if err != nil {
			log.Fatal(err)
		}
	}
}
```

And running instrumentation tests:

```go
func runTests(
	client adbaster.Client,
	device model.DeviceIdentifier,
	testPackage,
	runner,
	testClass string,
) {
	eventStream, instrumentationOutput, err := client.RunInstrumentationTests(
		device,
		model.InstrumentationParams{
			TestPackage: testPackage,
			Runner:      runner,
			Arguments: model.InstrumentationArguments{
				"testType":                            "INSTRUMENTATION",
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
		fmt.Println(event)
	}
}
```

For more details you can looks at [example code.](example/main.go)