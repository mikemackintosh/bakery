package pantry

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.sc-corp.net/Snapchat/ce-wimclasshero/cmd/ui"
)

// DownloadFile will download the source file (remote) to the dest (local) path
func DownloadFile(source, destination string, checksum interface{}) error {
	out, err := os.Create(destination)
	defer out.Close()

	// Make GET request
	resp, err := http.Get(source)
	if err != nil {
		return fmt.Errorf("Invalid download request, %s", err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		//ui.Warning("Error Downloading File: " + err.Error())
		return fmt.Errorf("Invalid server response")
	}

	// Get the file size
	fsize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// Create our progress reporter and pass it to be used alongside our writer
	counter := ui.NewWriteCounter(fsize)
	counter.Start()

	// Writer the body to file
	hash := sha256.New()
	_, err = io.Copy(io.MultiWriter(hash, out), io.TeeReader(resp.Body, counter))
	if err != nil {
		//ui.Warning()
		return fmt.Errorf("Error Downloading File: " + err.Error())
	}

	counter.Finish()

	fileChecksum := hex.EncodeToString(hash.Sum(nil))
	if checksum != nil {
		if checksum.(string) != fileChecksum {
			return fmt.Errorf("Failed to validate file. Want %s but have %s", checksum, fileChecksum)
		}
	}

	return nil
}
