package trapcmd

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
)

type InvocationResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func stream(reader io.Reader, writer io.Writer, doneChan chan<- bool, cancelChan <-chan bool) {
	b := make([]byte, 1)
	for {
		select {
		case <-cancelChan:
			doneChan <- true
			return
		default:
			inN, inErr := reader.Read(b)
			if inN > 0 {
				_, outErr := writer.Write(b)
				if outErr != nil {
					doneChan <- false
					return
				}
			}
			if inErr == io.EOF {
				doneChan <- true
				return
			}
			if inErr != nil {
				doneChan <- false
				return
			}
		}
	}
}

func RunWrappedCommand(trapCmd *cobra.Command, invocation []string) (*InvocationResult, error) {
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

	outDoneChannel := make(chan bool, 1)
	outCancelChannel := make(chan bool, 1)
	errDoneChannel := make(chan bool, 1)
	errCancelChannel := make(chan bool, 1)

	go stream(stdoutReader, trapCmd.OutOrStdout(), outDoneChannel, outCancelChannel)
	go stream(stderrReader, trapCmd.ErrOrStderr(), errDoneChannel, errCancelChannel)

	signalsChannel := make(chan os.Signal, 1)
	exitChannel := make(chan int, 1)
	signal.Notify(signalsChannel, os.Interrupt, os.Kill)

	go func() {
		success := true
		completed := 0
		for {
			select {
			case outSuccess := <-outDoneChannel:
				success = success && outSuccess
				completed += 1
			case errSuccess := <-errDoneChannel:
				success = success && errSuccess
				completed += 1
			case <-signalsChannel:
				outCancelChannel <- true
				errCancelChannel <- true
				exitChannel <- 130
				return
			default:
				if completed == 2 {
					if success {
						exitChannel <- 0
					} else {
						exitChannel <- 1
					}
					return
				}
			}
		}
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

	coordinatorExitCode := <-exitChannel
	if coordinatorExitCode != 0 {
		exitCode = coordinatorExitCode
	}

	return &InvocationResult{ExitCode: exitCode, Stdout: outBuilder.String(), Stderr: errBuilder.String()}, err
}
