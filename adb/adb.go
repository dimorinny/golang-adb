package adb

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	lineSeparator = "\r\n"
)

type Client struct {
	Config Config
}

func NewClient(config Config) Client {
	return Client{
		Config: config,
	}
}

func (c *Client) Connect() error {
	response, err := c.executeCommand(
		"connect",
		fmt.Sprintf(
			"%s:%d",
			c.Config.Host,
			c.Config.Port,
		),
	)
	if err != nil {
		return err
	}

	successPrefixes := []string{
		"connected to",
		"already connected",
	}
	for _, prefix := range successPrefixes {
		if strings.Contains(response, prefix) {
			return nil
		}
	}

	return errors.New("Error connecting to adb server. " + response)
}

func (c *Client) DeviceInfo() (*Device, error) {
	response, err := c.executeShellCommand("getprop")
	if err != nil {
		return nil, err
	}

	return newDeviceFromOutput(response, lineSeparator), nil
}

func (c *Client) executeShellCommand(arguments ...string) (string, error) {
	return c.executeCommand(
		append([]string{"shell"}, arguments...)...,
	)
}

func (c *Client) executeCommand(arguments ...string) (string, error) {
	output, err := exec.Command(
		c.Config.AdbPath,
		arguments...,
	).Output()

	if err != nil {
		return "", errors.New(
			fmt.Sprintf(
				"Some error while executing: %v.",
				arguments,
			),
		)
	}

	return string(output), nil
}
