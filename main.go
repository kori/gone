package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

// Command line arguments.
var (
	hostname = flag.String("h", "teknik.io",
		"Which host to use. See hosts.json for the available hosts.")
)

// Structs for dealing with unmarshalling the hosts.json file.
// hosts.json contains the information for the hosts.

type Host struct {
	Name      string `json:"name"`
	UploadURL string `json:"uploadurl"`
	ReturnURL string `json:"returnurl"`
}

type Hosts struct {
	Host []Host `json:"hosts"`
}

// Struct to unmarshal the response from the server.
type teknikResponse struct {
	Result struct {
		URL string
	}
}

// Struct to unmarshal the response from the server.
type pomfResponse struct {
	Files []struct {
		URL string
	}
}

// getHost returns a host with its relevant info.
func getHost() (*Host, error) {
	// Read hosts list
	path := os.Getenv("XDG_CONFIG_HOME") + "/gone/hosts.json"
	hostfile, err := ioutil.ReadFile(path)
	check(err)

	// Unmarshal hosts list and get the upload URL
	var hs Hosts
	json.Unmarshal(hostfile, &hs)
	check(err)
	for _, h := range hs.Host {
		if h.Name == *hostname {
			return &h, nil
		}
	}
	return nil, errors.New("getHost: host not in list")
}

// prepareUpload takes a path to a file and returns a request.
func (h *Host) prepareUpload(filepath string) *http.Request {
	// Read file to be uploaded.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(filepath)
	check(err)

	// Write to form
	fileForm, err := w.CreateFormFile("file", filepath)
	check(err)
	io.Copy(fileForm, f)
	// Close file and writer.
	f.Close()
	w.Close()

	req, err := http.NewRequest("POST", h.UploadURL, &b)
	check(err)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req
}

func (h *Host) upload(req *http.Request) string {
	// Prepare client and call API
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
				MaxVersion: tls.VersionTLS11,
			},
		},
	}

	// Perform request
	res, err := client.Do(req)
	check(err)

	// Unmarshal and return response.
	if *hostname == "teknik" {
		var r teknikResponse
		dec := json.NewDecoder(res.Body)
		dec.Decode(&r)

		return h.ReturnURL + r.Result.URL
	}

	var r pomfResponse
	dec := json.NewDecoder(res.Body)
	dec.Decode(&r)

	return h.ReturnURL + r.Files[0].URL
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	h, err := getHost()
	check(err)

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: gone -h [host] [files]")
		os.Exit(1)
	} else {
		// Upload each file.
		for _, p := range flag.Args() {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				if exists(file) {
					r := h.prepareUpload(file)
					fmt.Println(h.upload(r))
				} else {
					log.Fatal("file doesn't exist: ", file)
				}
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
