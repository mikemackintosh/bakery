package pantry

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/mikemackintosh/bakery/cli"
	"github.com/mikemackintosh/bakery/config"
	"github.com/zclconf/go-cty/cty"
)

type Shell struct {
	PantryItem
	Name      string   `hcl:"name,label"`
	Config    hcl.Body `hcl:",remain"`
	Script    string   `json:"script"`
	DependsOn []string `json:"depends_on"`
}

// identifies the DMG spec
var shellSpec = &hcldec.ObjectSpec{
	"depends_on": dependsOn,
	"script": &hcldec.AttrSpec{
		Name:     "script",
		Required: true,
		Type:     cty.String,
	},
}

//
func (d *Shell) Parse(evalContext *hcl.EvalContext) error {
	cli.Debug(cli.INFO, "Preparing Shell", d.Name)
	cfg, diags := hcldec.Decode(d.Config, shellSpec, evalContext)
	if len(diags) != 0 {
		for _, diag := range diags {
			cli.Debug(cli.INFO, "\t#", diag)
		}
		return fmt.Errorf("%s", diags.Errs()[0])
	}

	err := d.Populate(cfg, d)
	if err != nil {
		return err
	}

	return nil
}

func (d *Shell) Bake() {
	var tmpFile = config.Registry.TempDir + fmt.Sprintf("/%x.sh", sha256.Sum256([]byte(d.Script)))[:14]
	err := ioutil.WriteFile(tmpFile, []byte(d.Script), 0744)
	if err != nil {
		cli.Debug(cli.ERROR, fmt.Sprintf("Error writing script to %s", tmpFile), err)
	}

	cli.Debug(cli.INFO, fmt.Sprintf("Running script %s", tmpFile), err)

	o, err := RunCommand([]string{
		"/bin/bash",
		"-c",
		tmpFile})
	if err != nil {
		cli.Debug(cli.ERROR, fmt.Sprintf("Error running %s", tmpFile), err)
	}

	cli.Debug(cli.INFO, "\t->", o.String())
}

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
