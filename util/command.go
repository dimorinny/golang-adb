package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

func ExecuteCommandWithStreamOutput(command string, arguments ...string) (<-chan string, error) {
	cmd := exec.Command(
		command,
		arguments...,
	)

	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	reader := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(reader)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	result := make(chan string)
	go func() {
		for scanner.Scan() {
			result <- scanner.Text()
		}
		close(result)
	}()

	return result, nil
}

func ExecuteCommand(command string, arguments ...string) (string, error) {
	output, err := exec.Command(
		command,
		arguments...,
	).Output()

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
