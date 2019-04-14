package adb

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dimorinny/golang-adb/adb/instrumentation"
	"github.com/dimorinny/golang-adb/model"
	"github.com/dimorinny/golang-adb/util"
)

type Client struct {
	Config         Config
	PrintResponses bool

	instrumentationParser *instrumentation.Parser
}

func NewClient(config Config, printResponses bool) *Client {
	return &Client{
		Config:         config,
		PrintResponses: printResponses,

		instrumentationParser: instrumentation.NewParser(config.LineSeparator),
	}
}

func (c *Client) Devices() ([]model.DeviceIdentifier, error) {
	response, err := c.executeCommand(
		"devices",
	)

	if err != nil {
		return nil, err
	}

	return newDevicesIdentifiersFromOutput(response, c.Config.LineSeparator), nil
}

func (c *Client) DeviceInfo(device model.DeviceIdentifier) (*model.Device, error) {
	response, err := c.executeShellCommand(
		device,
		"getprop",
	)
	if err != nil {
		return nil, err
	}

	return newDeviceFromOutput(response, c.Config.LineSeparator), nil
}

func (c *Client) Push(device model.DeviceIdentifier, from, to string) error {
	_, err := c.executeDeviceCommand(
		device,
		"push",
		from,
		to,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Pull(device model.DeviceIdentifier, from, to string) error {
	_, err := c.executeDeviceCommand(
		device,
		"pull",
		from,
		to,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Install(device model.DeviceIdentifier, from string) error {
	if _, err := os.Stat(from); os.IsNotExist(err) {
		return err
	}

	to := fmt.Sprintf("/data/local/tmp/%s", path.Base(from))
	err := c.Push(device, from, to)
	if err != nil {
		return err
	}

	defer func() {
		c.executeShellCommand(device, "rm", to)
	}()

	response, err := c.executeShellCommand(
		device,
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

func (c *Client) RunInstrumentationTests(
	device model.DeviceIdentifier,
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

	output, err := c.executeShellStreamCommand(device, arguments...)
	if err != nil {
		return nil, nil, err
	}

	eventStream, instrumentationOutputStream := c.instrumentationParser.Process(output)

	return eventStream, instrumentationOutputStream, nil
}

func (c *Client) ClearApplicationData(device model.DeviceIdentifier, applicationPackage string) error {
	output, err := c.executeShellCommand(
		device,
		"pm",
		"clear",
		applicationPackage,
	)
	if err != nil {
		return err
	}

	return detectErrorInClearApplicationDataOutput(output)
}

func (c *Client) Logcat(device model.DeviceIdentifier) (<-chan string, error) {
	outputStream, err := c.executeDeviceStreamCommand(device, "logcat")
	if err != nil {
		return nil, err
	}

	return outputStream, nil
}

func (c *Client) printResponseForCommand(command, response string) {
	fmt.Println(fmt.Sprintf("Output for command - %s:", command))
	fmt.Println(response)
}

// shell command utils
func (c *Client) executeShellCommand(device model.DeviceIdentifier, arguments ...string) (string, error) {
	return c.executeDeviceCommand(
		device,
		append(
			[]string{"shell"},
			arguments...,
		)...,
	)
}

func (c *Client) executeDeviceCommand(device model.DeviceIdentifier, arguments ...string) (string, error) {
	return c.executeCommand(
		append(
			[]string{"-s", string(device)},
			arguments...,
		)...,
	)
}

func (c *Client) executeCommand(arguments ...string) (string, error) {
	output, err := exec.Command(
		c.Config.AdbPath,
		arguments...,
	).Output()

	if c.PrintResponses {
		c.printResponseForCommand(
			strings.Join(arguments, " "),
			string(output),
		)
	}

	if err != nil {
		return "", errors.New(
			fmt.Sprintf(
				"some error while executing: %v. output: %s",
				arguments,
				err,
			),
		)
	}

	return string(output), nil
}

// shell stream command utils
func (c *Client) executeShellStreamCommand(device model.DeviceIdentifier, arguments ...string) (<-chan string, error) {
	return c.executeDeviceStreamCommand(
		device,
		append(
			[]string{"shell"},
			arguments...,
		)...,
	)
}

func (c *Client) executeDeviceStreamCommand(device model.DeviceIdentifier, arguments ...string) (<-chan string, error) {
	return c.executeStreamCommand(
		append(
			[]string{"-s", string(device)},
			arguments...,
		)...,
	)
}

func (c *Client) executeStreamCommand(arguments ...string) (<-chan string, error) {
	return util.ExecCommandWithStreamOutput(
		c.Config.AdbPath,
		arguments...,
	)
}
