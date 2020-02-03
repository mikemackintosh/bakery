package pantry

import (
	"testing"
	"net/http"
)

var testDownloadFile = []struct {
	Source      string
	Destination string
	Checksum    string
}{
	{
		Source:      "http://localhost:3000/test_download.txt",
		Destination: "/tmp/test_download.txt",
		Checksum:    "241fc234ac0c3cd1d85cd766800a72dc79a46f62fc29e99b1b3e2a41bfeece2b",
	},
}

func TestDownloadFile(t *testing.T) {
		http.Handle("/test_download.txt", http.FileServer(http.Dir("../testing/fixtures/www")))
		go http.ListenAndServe("localhost:3000", nil)

		for _, test := range testDownloadFile {
			err := DownloadFile(test.Source, test.Destination, test.Checksum)
			if err != nil {
				t.Errorf("DownloadFile failed with %s", err)
			}
		}

}
