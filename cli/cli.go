package cli

import (
	"flag"
	"fmt"
	"os"
)

// flag options for CLI
var (
	FlagConfig    string
	FlagRecipe    string
	FlagDebug     bool
	FlagVerbosity int

	// severityName maps severity const's to string names
	severityName = []string{
		INFO:    "INFO",
		ERROR:   "ERROR",
		WARNING: "WARNING",
		DEBUG:   "DEBUG",
		DEBUG2:  "DEBUG2",
		DEBUG3:  "DEBUG3",
	}
)

// iota is used to set incrementing value for constant list
const (
	INFO = iota
	ERROR
	WARNING
	DEBUG
	DEBUG2
	DEBUG3
)

// Init flags
func init() {
	flag.StringVar(&FlagConfig, "c", "manifest.yml", "Configuration file")
	flag.StringVar(&FlagRecipe, "r", "config.yum", "Client recipe file")
	flag.BoolVar(&FlagDebug, "d", false, "When enabled, turns on debugging")
	flag.IntVar(&FlagVerbosity, "v", 1, "Sets output verbosity level")
}

// Debug prints out debug messaging with log levels when set
func Debug(verbosity int, msg string, value interface{}) string {
	var output string
	if FlagDebug && verbosity <= FlagVerbosity {
		output = fmt.Sprintf("[%6s ] ", severityName[verbosity])
		if value != nil {
			output = fmt.Sprintf("%s%s: %v\n", output, msg, value)
			fmt.Print(output)
			return output
		}

		output = fmt.Sprintf("%s%s\n", output, msg)
	}

	fmt.Print(output)
	return output
}

// ErrorAndExit will print an error and exit the program
func ErrorAndExit(err error) {
	fmt.Printf("%s", err)
	os.Exit(1)
}
