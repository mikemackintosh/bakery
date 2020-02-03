package cli

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	pb "gopkg.in/cheggaaa/pb.v1"
)

// Header formats and pretty prints a header
func Header(str string) {
	color.New(color.FgCyan).Add(color.Bold).Println(str)
}

// SubHeader formats and pretty prints a sub-header
func SubHeader(str, str2 interface{}) {
	var left = color.New(color.FgYellow).Add(color.Bold)
	var right = color.New(color.FgMagenta)
	left.Printf("%v: ", str)
	right.Printf("%v\n", str2)
}

// Warning bolds and continues
func Warning(err string) {
	var warning = color.New(color.FgYellow).Add(color.Bold)
	warning.Print("\nWarning! ")
	warning.DisableColor()
	fmt.Printf("%s\n\n", err)
}

// Error bolds and exits
func Error(err string) {
	var warning = color.New(color.FgRed).Add(color.Bold)
	warning.Println("\n!! Error Encountered !!")
	warning.DisableColor()
	fmt.Printf("\t%s\n\n", err)
	os.Exit(1)
}

// Success bolds and prints green
func Success(msg string) {
	var success = color.New(color.FgGreen).Add(color.Bold)
	success.Printf("%v\n", msg)
}

// PrintKV prints out a key: value with columnar highlightins
func PrintKV(key interface{}, value string) {
	var printKey = color.New(color.FgWhite)
	var printValue = color.New(color.FgGreen)

	printKey.Printf("%v: ", key)
	printValue.Printf("%s\n", value)
}

// PrintErrorKV prints out a key: value with columnar highlightins
func PrintErrorKV(key interface{}, value interface{}) {
	var printKey = color.New(color.FgWhite)
	var printValue = color.New(color.FgRed)

	printKey.Printf("%v: ", key)
	printValue.Printf("%v\n", value)
}

// SprintKV prints out a key: value with columnar highlightins
func SprintKV(key interface{}, value string) string {
	var printKey = color.New(color.FgWhite)
	var printValue = color.New(color.FgGreen)

	return printKey.Sprintf("%v ", key) + printValue.Sprintf("%s\n", value)
}

// SprintErrorKV prints out a key: value with columnar highlightins
func SprintErrorKV(key interface{}, value string) string {
	var printKey = color.New(color.FgWhite)
	var printValue = color.New(color.FgRed)

	return printKey.Sprintf("%v ", key) + printValue.Sprintf("%s\n", value)
}

// Ask will ask a question and wait for input
func Ask(question string) string {
	reader := bufio.NewReader(os.Stdin)
	var askPrint = color.New(color.FgGreen)
	askPrint.Printf(question + ": ")
	text, _ := reader.ReadString('\n')
	return text
}

const RefreshRate = time.Millisecond * 100

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	n   int // bytes read so far
	bar *pb.ProgressBar
}

func NewWriteCounter(total int) *WriteCounter {
	b := pb.New(total)
	b.SetRefreshRate(RefreshRate)
	b.ShowTimeLeft = true
	b.ShowSpeed = true
	b.SetUnits(pb.U_BYTES)

	return &WriteCounter{
		bar: b,
	}
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	wc.n += len(p)
	wc.bar.Set(wc.n)
	return wc.n, nil
}

func (wc *WriteCounter) Start() {
	wc.bar.Start()
}

func (wc *WriteCounter) Finish() {
	wc.bar.Finish()
}
