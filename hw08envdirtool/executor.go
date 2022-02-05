package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmdArgs []string, env Environment) (returnCode int, err error) {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	envStr := make([]string, 0, len(env))
	for key, val := range env {
		envStr = append(envStr, fmt.Sprintf("%s=%s", key, val))
	}

	cmd.Env = append(os.Environ(), envStr...)

	if err := cmd.Run(); err != nil {
		var exitError *exec.ExitError

		if errors.As(err, &exitError) {
			return exitError.ExitCode(), err
		}

		return 0, err
	}

	return 0, nil
}
