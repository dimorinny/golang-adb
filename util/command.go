package util

import (
	"bufio"
	"io"
	"os/exec"
)

func ExecCommandWithStreamOutput(command string, arguments ...string) (<-chan string, error) {
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
