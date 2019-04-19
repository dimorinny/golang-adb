package adb

import (
	"errors"
	"fmt"
	"github.com/dimorinny/golang-adb"
	"github.com/dimorinny/golang-adb/adb/instrumentation"
	"github.com/dimorinny/golang-adb/model"
	"github.com/dimorinny/golang-adb/util"
	"os"
	"path"
	"strings"
)

type Device struct {
	Identifier model.DeviceIdentifier

	adbPath        string
	printResponses bool

	instrumentationParser *instrumentation.Parser
}

func newDevice(
	identifier model.DeviceIdentifier,
	adbPath string,
	printResponses bool,
) golangadb.Device {
	return &Device{
		Identifier: identifier,

		adbPath:        adbPath,
		printResponses: printResponses,

		instrumentationParser: instrumentation.NewParser(lineSeparator),
	}
}

func (d *Device) Info() (*model.Device, error) {
	response, err := d.executeShellCommand("getprop")
	if err != nil {
		return nil, err
	}

	return newDeviceFromOutput(d.Identifier, response, lineSeparator), nil
}

func (d *Device) Push(from, to string) error {
	_, err := d.executeCommand(
		"push",
		from,
		to,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) Pull(from, to string) error {
	_, err := d.executeCommand(
		"pull",
		from,
		to,
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) Install(from string) error {
	if _, err := os.Stat(from); os.IsNotExist(err) {
		return err
	}

	to := fmt.Sprintf("/data/local/tmp/%s", path.Base(from))
	err := d.Push(from, to)
	if err != nil {
		return err
	}

	defer func() {
		output, err := d.executeShellCommand(
			"rm",
			to,
		)
		if err != nil {
			fmt.Println(
				fmt.Sprintf("Failed to remove application from device. Output: %s", output),
			)
		}
	}()

	response, err := d.executeShellCommand(
		"pm",
		"install",
		"-r",
		to,
	)

	if err != nil {
		return err
	}

	if !strings.HasSuffix(strings.TrimSpace(response), "Success") {
		return errors.New(
			fmt.Sprintf(
				"application installign failure with output: %s",
				response,
			),
		)
	}

	return nil
}

func (d *Device) RunInstrumentationTests(
	params model.InstrumentationParams,
) (<-chan instrumentation.Event, <-chan string, error) {
	if params.TestPackage == "" || params.Runner == "" {
		return nil, nil, errors.New(
			"test package and runner params is required in RunInstrumentationTests method",
		)
	}

	var arguments []string
	arguments = append(
		arguments,
		"am",
		"instrument",
		"-w",
		"-r",
	)

	instrumentationArguments := util.Join(
		"-e %s %s",
		params.Arguments,
	)

	if len(instrumentationArguments) > 0 {
		arguments = append(
			arguments,
			instrumentationArguments...,
		)
	}

	arguments = append(
		arguments,
		fmt.Sprintf(
			"%s/%s",
			params.TestPackage,
			params.Runner,
		),
	)

	output, err := d.executeShellStreamCommand(arguments...)
	if err != nil {
		return nil, nil, err
	}

	eventStream, instrumentationOutputStream := d.instrumentationParser.Process(output)

	return eventStream, instrumentationOutputStream, nil
}

func (d *Device) ClearApplicationData(applicationPackage string) error {
	output, err := d.executeShellCommand(
		"pm",
		"clear",
		applicationPackage,
	)
	if err != nil {
		return err
	}

	return detectErrorInClearApplicationDataOutput(output)
}

func (d *Device) Logcat() (<-chan string, error) {
	outputStream, err := d.executeStreamCommand("logcat")
	if err != nil {
		return nil, err
	}

	return outputStream, nil
}

func (d *Device) printResponseForCommand(command, response string) {
	fmt.Println(fmt.Sprintf("Output for command - %s:", command))
	fmt.Println(response)
}

func (d *Device) executeCommand(arguments ...string) (string, error) {
	return d.executeAdbCommand(
		append(
			[]string{"-s", string(d.Identifier)},
			arguments...,
		)...,
	)
}

func (d *Device) executeShellCommand(arguments ...string) (string, error) {
	return d.executeCommand(
		append(
			[]string{"shell"},
			arguments...,
		)...,
	)
}

func (d *Device) executeAdbCommand(arguments ...string) (string, error) {
	output, err := util.ExecuteCommand(
		d.adbPath,
		arguments...,
	)

	if d.printResponses {
		d.printResponseForCommand(
			strings.Join(arguments, " "),
			string(output),
		)
	}

	return output, err
}

// shell stream command utils
func (d *Device) executeStreamCommand(arguments ...string) (<-chan string, error) {
	return d.executeAdbStreamCommand(
		append(
			[]string{"-s", string(d.Identifier)},
			arguments...,
		)...,
	)
}

func (d *Device) executeShellStreamCommand(arguments ...string) (<-chan string, error) {
	return d.executeStreamCommand(
		append(
			[]string{"shell"},
			arguments...,
		)...,
	)
}

func (d *Device) executeAdbStreamCommand(arguments ...string) (<-chan string, error) {
	return util.ExecuteCommandWithStreamOutput(
		d.adbPath,
		arguments...,
	)
}
