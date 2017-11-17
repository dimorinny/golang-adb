package adb

import (
	"errors"
	"fmt"
	"github.com/dimorinny/adbaster/model"
	"github.com/dimorinny/adbaster/util"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	lineSeparator = "\n"
)

type Client struct {
	Config         Config
	PrintResponses bool
}

func NewClient(config Config, printResponses bool) Client {
	return Client{
		Config:         config,
		PrintResponses: printResponses,
	}
}

func (c *Client) Devices() ([]model.DeviceIdentifier, error) {
	response, err := c.executeCommand(
		"devices",
	)

	if err != nil {
		return nil, err
	}

	return newDevicesIdentifiersFromOutput(response, lineSeparator), nil
}

func (c *Client) DeviceInfo(device model.DeviceIdentifier) (*model.Device, error) {
	response, err := c.executeShellCommand(
		device,
		"getprop",
	)
	if err != nil {
		return nil, err
	}

	return newDeviceFromOutput(response, lineSeparator), nil
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
) (*model.InstrumentationResult, error) {
	if params.Package == "" || params.Runner == "" {
		return nil, errors.New(
			"package and runner params is required in RunInstrumentationTests method",
		)
	}

	arguments := []string{}
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
			params.Package,
			params.Runner,
		),
	)

	response, err := c.executeShellCommand(device, arguments...)
	if err != nil {
		return nil, err
	}

	if strings.Contains(response, "INSTRUMENTATION_STATUS: Error") {
		return nil, errors.New(
			fmt.Sprintf(
				"error while starting instrumental tests with output: %s",
				response,
			),
		)
	}

	return newInstrumentationResultFromOutput(response), nil
}

func (c *Client) printResponseForCommand(command, response string) {
	fmt.Println(fmt.Sprintf("Output for command - %s:", command))
	fmt.Println(response)
}

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
