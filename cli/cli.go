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
	FlagTempDir   string
	FlagBundle    bool
	FlagDebug     bool
	FlagVerbosity int

	// severityName maps severity const's to string names
	severityName = []Severity{
		INFO:    Severity{Name: "INFO", Color: "\033[38;5;45m"},
		ERROR:   Severity{Name: "ERROR", Color: "\033[38;5;196"},
		WARNING: Severity{Name: "WARNING", Color: "\033[38;5;214m"},
		DEBUG:   Severity{Name: "DEBUG", Color: "\033[38;5;45m"},
		DEBUG2:  Severity{Name: "DEBUG2", Color: "\033[38;5;45m"},
		DEBUG3:  Severity{Name: "DEBUG3", Color: "\033[38;5;45m"},
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

type Severity struct {
	Name  string
	Color string
}

// Init flags
func init() {
	flag.StringVar(&FlagConfig, "c", "manifest.yml", "Configuration file")
	flag.StringVar(&FlagRecipe, "r", "config.yum", "Client recipe file")
	flag.StringVar(&FlagTempDir, "-temp-dir", " /var/bakery/tmp", "Temporary resource directory")
	flag.BoolVar(&FlagBundle, "b", false, "Bundle client config with binary")
	flag.BoolVar(&FlagDebug, "d", false, "When enabled, turns on debugging")
	flag.IntVar(&FlagVerbosity, "v", 1, "Sets output verbosity level")
}

// Debug prints out debug messaging with log levels when set
func Debug(verbosity int, msg string, value interface{}) string {
	var output string
	if FlagDebug && verbosity <= FlagVerbosity {
		output = fmt.Sprintf("\033[0m[%7s ] %s", severityName[verbosity].Name, severityName[verbosity].Color)
		if value != nil {
			output = fmt.Sprintf("%s%s: %v\n", output, msg, value)
			fmt.Print(output)
			return output
		}

		output = fmt.Sprintf("%s%s\033[0m\n", output, msg)
	}

	fmt.Print(output)
	return output
}

// ErrorAndExit will print an error and exit the program
func ErrorAndExit(err error) {
	fmt.Printf("%s", err)
	os.Exit(1)
}
