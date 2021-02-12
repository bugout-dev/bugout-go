package trapcmd

import (
	"bytes"
	"os/exec"
)

type InvocationResult struct {
	ExitCode  int
	OutBuffer *bytes.Buffer
	ErrBuffer *bytes.Buffer
}

func RunWrappedCommand(invocation []string) (InvocationResult, error) {
	cmd := exec.Command(invocation[0], invocation[1:]...)
	var outBuffer, errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	exitCode := 0

	var err error
	runErr := cmd.Run()
	if runErr != nil {
		if cmdErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = cmdErr.ExitCode()
		} else {
			// runWrappedCommand only returns a non-nil error if it received an error from cmd.Run
			// wasn't an exit error for the command (i.e. wasn't a non-zero exit code for the
			// invocation).
			err = runErr
		}
	}

	return InvocationResult{ExitCode: exitCode, OutBuffer: &outBuffer, ErrBuffer: &errBuffer}, err
}
