package pantry

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testDownloadFile = []struct {
	Source      string
	Destination string
	Checksum    string
}{
	{
		Source:      "/test_download.txt",
		Destination: "/tmp/test_download.txt",
		Checksum:    "241fc234ac0c3cd1d85cd766800a72dc79a46f62fc29e99b1b3e2a41bfeece2b",
	},
}

func TestDownloadFile(t *testing.T) {
	var filename = "/test_download.txt"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == filename {
			b, _ := ioutil.ReadFile("../../testing/fixtures/www/test_download.txt")
			_, _ = w.Write(b)
			return
		}
	}))

	for _, test := range testDownloadFile {
		err := DownloadFile(ts.Listener.Addr().String()+test.Source, test.Destination, test.Checksum)
		if err != nil {
			t.Errorf("DownloadFile failed with %s", err)
		}
	}

}
