package cli

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

/*
 * Test_debug tests debug messaging output and leveling
 */
var debugOutputTest = []struct {
	Verbosity int
	Debug     bool
	Message   string
	Expected  string
}{
	{
		ERROR, true, "Test", "\033[0m[  ERROR ] \033[38;5;196mTest\033[0m\n",
	},
	{
		ERROR, false, "Test", "",
	},
	{
		DEBUG, true, "Debug Message", "\033[0m[  DEBUG ] \033[38;5;45mDebug Message\033[0m\n",
	},
}

func TestDebug(t *testing.T) {
	for _, test := range debugOutputTest {
		FlagVerbosity = test.Verbosity
		FlagDebug = test.Debug
		output := Debug(test.Verbosity, test.Message, nil)
		if output != test.Expected {
			t.Errorf("wanted '%s' but got '%s'", test.Expected, output)
		}
	}
}

/*
 * TestErrorAndExit tests that the function call performs an OS exit
 */
func TestErrorAndExit(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		ErrorAndExit(fmt.Errorf("Error"))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestErrorAndExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
