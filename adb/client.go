package adb

import (
	"fmt"
	"github.com/dimorinny/golang-adb"
	"github.com/dimorinny/golang-adb/util"
	"strings"

	"github.com/dimorinny/golang-adb/adb/instrumentation"
)

type Client struct {
	adbPath        string
	printResponses bool

	instrumentationParser *instrumentation.Parser
}

func NewClient(adbPath string, printResponses bool) golangadb.Client {
	return &Client{
		adbPath:        adbPath,
		printResponses: printResponses,

		instrumentationParser: instrumentation.NewParser(lineSeparator),
	}
}

func (c *Client) Devices() ([]golangadb.Device, error) {
	response, err := c.executeAdbCommand(
		"devices",
	)

	if err != nil {
		return nil, err
	}

	var result []golangadb.Device

	devices := newDevicesIdentifiersFromOutput(response, lineSeparator)

	for _, identifier := range devices {
		result = append(
			result,
			newDevice(
				identifier,
				c.adbPath,
				c.printResponses,
			),
		)
	}

	return result, nil
}

func (c *Client) printResponseForCommand(command, response string) {
	fmt.Println(fmt.Sprintf("Output for command - %s:", command))
	fmt.Println(response)
}

func (c *Client) executeAdbCommand(arguments ...string) (string, error) {
	output, err := util.ExecuteCommand(
		c.adbPath,
		arguments...,
	)

	if c.printResponses {
		c.printResponseForCommand(
			strings.Join(arguments, " "),
			string(output),
		)
	}

	return output, err
}
