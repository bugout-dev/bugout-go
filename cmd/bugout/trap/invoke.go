package trapcmd

import (
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type InvocationResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func RunWrappedCommand(trapCmd *cobra.Command, invocation []string) (*InvocationResult, error) {
	var wg sync.WaitGroup
	cmd := exec.Command(invocation[0], invocation[1:]...)
	cmd.Stdin = trapCmd.InOrStdin()

	outReader, stdoutPipeErr := cmd.StdoutPipe()
	if stdoutPipeErr != nil {
		return &InvocationResult{}, stdoutPipeErr
	}
	errReader, stderrPipeErr := cmd.StderrPipe()
	if stderrPipeErr != nil {
		return &InvocationResult{}, stderrPipeErr
	}

	var outBuilder, errBuilder strings.Builder
	stdoutReader := io.TeeReader(outReader, &outBuilder)
	stderrReader := io.TeeReader(errReader, &errBuilder)

	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(trapCmd.OutOrStdout(), stdoutReader)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(trapCmd.ErrOrStderr(), stderrReader)
	}()

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
	wg.Wait()
	return &InvocationResult{ExitCode: exitCode, Stdout: outBuilder.String(), Stderr: errBuilder.String()}, err
}
