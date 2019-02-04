package pantry

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

// TestRunCommandOutput used to evaluate a successful test
var TestRunCommandOutput = `Result
FOO:BAR`

type CommandResponse struct {
	Command *exec.Cmd
	Raw     string
	Error   string
}

// StreamCommand will stream the output of the command to the specified buffer with
// or without a prefix
func StreamCommand(prefix string, cmdArgs []string) error {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()

	prefix = "\t"

	go func() {
		_, errStdout = copyOutputAndCapture(color.Output, stdoutIn, prefix, color.FgGreen)
	}()

	go func() {
		_, errStderr = copyOutputAndCapture(color.Output, stderrIn, prefix, color.FgRed)
	}()

	err := cmd.Wait()
	if err != nil {
		return fmt.Errorf("Failed to run command with error %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		return fmt.Errorf("No output received\n")
	}

	fmt.Println()
	return nil
}

func copyOutputAndCapture(w io.Writer, r io.Reader, prefix string, c color.Attribute) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			comp := color.New(color.FgBlue)
			comp.Fprintf(w, "%s", prefix)

			output := color.New(c)

			out = append(out, d...)
			_, err := output.Fprint(w, strings.TrimSpace(string(d)))
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

// RunCommand returns true if the audit passed, or command was successful
func RunCommand(cmdArgs []string) (*CommandResponse, error) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		return &CommandResponse{
			Command: nil,
			Raw:     TestRunCommandOutput,
		}, nil
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	res, err := cmd.Output()
	return &CommandResponse{
		Command: cmd,
		Raw:     strings.TrimSpace(string(res)),
	}, err
}

func (r *CommandResponse) String() string {
	return r.Raw
}

// ByLine will split the command output by line
func (r *CommandResponse) ByLine() []string {
	var output []string
	response := r.String()
	if runtime.GOOS == "windows" {
		output = strings.Split(response, "\r\r\n")
	} else {
		output = strings.Split(response, "\n")
	}

	// Trim the output of unneccessary whitespace in the output
	for k, v := range output {
		output[k] = strings.TrimSpace(v)
	}

	return output
}

// Grep returns
func (r *CommandResponse) Grep(grep string) []string {
	var matches []string
	for _, line := range r.ByLine() {
		if strings.Contains(line, grep) {
			matches = append(matches, line)
		}
	}

	return matches
}

// SplitColon returns a split
func (r *CommandResponse) SplitColon() []string {
	return strings.Split(r.String(), ":")
}
