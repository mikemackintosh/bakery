package pantry

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/mikemackintosh/bakery/cli"
)

// DownloadFile will download the source file (remote) to the dest (local) path
func DownloadFile(source, destination string, checksum interface{}) error {
	if FileExists(destination) {
		cli.Debug(cli.INFO, fmt.Sprintf("\t-> Destination file %s already exists", destination), nil)
		if checksum != nil {
			f, err := os.Open(destination)
			if err != nil {
				return err
			}
			defer f.Close()
			hash := sha256.New()
			if _, err := io.Copy(hash, f); err != nil {
				return err
			}

			fileHash := hex.EncodeToString(hash.Sum(nil))
			if fileHash != checksum.(string) {
				cli.Debug(cli.INFO, fmt.Sprintf("\t-> File with hash %s detected, but want %s, removing...", fileHash, checksum.(string)), nil)
				err := os.RemoveAll(destination)
				if err != nil {
					return fmt.Errorf("Error removing invalidated file: %s", err)
				}
			} else {
				cli.Debug(cli.INFO, "\t-> Using existing destination file", nil)
				return nil
			}
		}
	}

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
		return fmt.Errorf("Invalid server response")
	}

	// Get the file size
	fsize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// Create our progress reporter and pass it to be used alongside our writer
	counter := cli.NewWriteCounter(fsize)
	counter.Start()

	// Writer the body to file
	hash := sha256.New()
	_, err = io.Copy(io.MultiWriter(hash, out), io.TeeReader(resp.Body, counter))
	if err != nil {
		//cli.Warning()
		return fmt.Errorf("Error Downloading File: " + err.Error())
	}

	counter.Finish()

	fileChecksum := hex.EncodeToString(hash.Sum(nil))
	if len(checksum.(string)) > 0 {
		if checksum.(string) != fileChecksum {
			return fmt.Errorf("Failed to validate file. Want %s but have %s", checksum.(string), fileChecksum)
		}
	}

	return nil
}
