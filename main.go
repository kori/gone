package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

func Upload(url string, filepath string) (string, error) {
	// Read file to be uploaded.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	// Write to form
	fileForm, err := w.CreateFormFile("file", filepath)
	if err != nil {
		return "", err
	}
	io.Copy(fileForm, f)
	// Close file and writer.
	f.Close()
	w.Close()

	req, err := http.NewRequest("POST", url, &b)
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
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var r struct {
		Result struct {
			URL string
		}
	}

	dec := json.NewDecoder(res.Body)
	dec.Decode(&r)

	return r.Result.URL, nil
}

func main() {
	var wg sync.WaitGroup

	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatalln("Usage: gone [files]")
	} else {
		// Upload each file.
		for _, p := range flag.Args() {
			wg.Add(1)
			go func(file string) {
				if exists(file) {
					link, err := Upload("https://api.teknik.io/v1/Upload", file)
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Println(file+":", link)

				} else {
					log.Fatal("file doesn't exist: ", file)
				}
				wg.Done()
			}(p)
		}
		wg.Wait()
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
