package pantry

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("GO_WANT_HELPER_PROCESS", "1")
}

/*
The followign test will check to make sure the RunCommand returned the expected output
*/
var testRunCommand = struct {
	Expected string
	Error    error
}{
	Expected: TestRunCommandOutput,
	Error:    nil,
}

func TestRunCommand(t *testing.T) {
	output, err := RunCommand([]string{"test"})
	if testRunCommand.Error != err {
		t.Fatalf("want %s but got %s", testRunCommand.Error, err)
	}

	if testRunCommand.Expected != output.Raw {
		t.Fatalf("want %s but got %s", testRunCommand.Expected, output.Raw)
	}
}

/*
The following test validated the successful implementation of the Grep functionality
*/
var testCommandResponseGrep = struct {
	Expected []string
	Error    error
}{
	Expected: []string{"FOO:BAR"},
	Error:    nil,
}

func TestCommandResponseGrep(t *testing.T) {
	output, err := RunCommand([]string{"test"})
	if testCommandResponseGrep.Error != err {
		t.Fatalf("want %s but got %s", testCommandResponseGrep.Error, err)
	}

	grep := output.Grep("FOO")
	if len(testCommandResponseGrep.Expected) != len(grep) {
		t.Fatalf("want %s but got %s", testCommandResponseGrep.Expected, grep)
	}

	if testCommandResponseGrep.Expected[0] != grep[0] {
		t.Fatalf("want %s but got %s", testCommandResponseGrep.Expected, grep)
	}
}

/*
The followign test will check to make sure the split colon funtionality is working as expected
*/
var testCommandResponseSplitColo = struct {
	Expected []string
	Error    error
}{
	Expected: []string{"FOO", "BAR"},
	Error:    nil,
}

func TestCommandResponseSplitColon(t *testing.T) {
	output, err := RunCommand([]string{"test"})
	if testCommandResponseSplitColo.Error != err {
		t.Fatalf("want %s but got %s", testCommandResponseSplitColo.Error, err)
	}

	grep := output.SplitColon()
	if len(testCommandResponseSplitColo.Expected) != len(grep) {
		t.Fatalf("want %s but got %s", testCommandResponseSplitColo.Expected, grep)
	}

	if testCommandResponseSplitColo.Expected[1] != grep[1] {
		t.Fatalf("want %s but got %s", testCommandResponseSplitColo.Expected, grep)
	}
}
