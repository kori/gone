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

func Upload(filepath string) (string, error) {
	// Read file to be uploaded.
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fi, err := os.Stat(filepath)
	if err != nil {
		return "", err
	}

	// Write to form
	form, err := w.CreateFormFile("file", fi.Name())
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(form, f); err != nil {
		return "", err
	}
	// Close file and writer.
	f.Close()
	w.Close()

	req, err := http.NewRequest("POST", "https://api.teknik.io/v1/Upload", &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	// Prepare client and call API
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
				MaxVersion: tls.VersionTLS11,
			},
		},
	}
	// Unmarshal and return response.
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var r struct {
		Result struct {
			URL string
		}
	}

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&r)

	resp.Body.Close()

	return r.Result.URL, nil
}
