package teknik

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const apiURL = "https://api.teknik.io/upload/post"

type ResponseFile struct {
	Results struct {
		File struct {
			URL string `json:"url"` // The direct URL of the uploaded file.
		} `json:"file"`
	} `json:"results"`
}

type ResponsePaste struct {
	Results struct {
		Pastes struct {
			URL string `json:"url"` // The direct url to the paste.
		} `json:"paste"`
	} `json:"results"`
}

func Upload(file string) string {
	// Set up a Transport with Teknik's TLS version info.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS11,
			MaxVersion: tls.VersionTLS11,
		},
	}

	// Read file to be uploaded.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	f, err := os.Open(file)
	check(err)

	// Write to form
	fileForm, err := w.CreateFormFile("file", file)
	check(err)

	io.Copy(fileForm, f)
	f.Close()

	err = w.Close()
	check(err)

	req, err := http.NewRequest("POST", apiURL, &b)
	check(err)

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{
		Transport: tr,
	}

	// Call API
	res, err := client.Do(req)
	check(err)

	// Print results after unmarshalling
	buf := make([]byte, res.ContentLength)
	io.ReadFull(res.Body, buf)

	var r []ResponseFile
	check(json.Unmarshal(buf, &r))

	return r[0].Results.File.URL
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
