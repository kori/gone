package teknik

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const apiURL = "https://api.teknik.io/upload/post"

type response struct {
	Results struct {
		File struct {
			URL string `json:"url"` // The direct URL to the uploaded file.
		} `json:"file"`
	} `json:"results"`
}

// Function Upload takes a path to a file and returns the URL from Teknik.
func Upload(file string) string {
	// Read file to be uploaded.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(file)
	check(err)

	// Write to form
	fileForm, err := w.CreateFormFile("file", file)
	check(err)
	io.Copy(fileForm, f)
	// Close file and writer.
	f.Close()
	w.Close()

	req, err := http.NewRequest("POST", apiURL, &b)
	check(err)
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Set up a Transport with Teknik's TLS version info, and then set up
	// the HTTP client.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS11,
			MaxVersion: tls.VersionTLS11,
		},
	}
	client := &http.Client{
		Transport: tr,
	}

	// Call API
	res, err := client.Do(req)
	check(err)

	// Return results after unmarshal.
	buf := make([]byte, res.ContentLength)
	io.ReadFull(res.Body, buf)
	var r []response
	check(json.Unmarshal(buf, &r))

	return r[0].Results.File.URL
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
