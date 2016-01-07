package pomf

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// Structs for dealing with unmarshalling the hosts.json file.
// hosts.json contains the information for the pomf-related hosts.
type host struct {
	Name      string
	UploadURL string
	ReturnURL string
}
type hostContainer struct {
	Hosts []host
}

// Struct to unmarshal the response from the server.
type response struct {
	Files []struct {
		URL string
	}
}

func getHost(hostname string) (*host, error) {
	// Read hosts list
	path := os.Getenv("XDG_CONFIG_HOME") + "/gone/hosts.json"
	hostfile, err := ioutil.ReadFile(path)
	check(err)

	// Unmarshal hosts list and get the upload URL
	var container hostContainer
	json.Unmarshal(hostfile, &container)
	check(err)
	for _, host := range container.Hosts {
		if host.Name == hostname {
			return &host, nil
		}
	}
	return nil, errors.New("getHost: host not in list")
}

// Function Upload takes the path for a file, and the host name for a
// pomf-based website, if it matches a host on the current host lists,
// it performs the upload to the server and returns the URL.
func Upload(file string, hostname string) string {
	// Read file to be uploaded.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(file)
	check(err)

	// Write to form
	fileForm, err := w.CreateFormFile("files[]", file)
	check(err)
	io.Copy(fileForm, f)
	// Close file and writer.
	f.Close()
	w.Close()

	// Get host's information.
	h, err := getHost(hostname)
	check(err)

	// Call API
	req, err := http.NewRequest("POST", h.UploadURL, &b)
	check(err)

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	check(err)

	// Unmarshal response.
	var r response
	dec := json.NewDecoder(res.Body)
	dec.Decode(&r)

	if h.Name == "mixtape.moe" || h.Name == "pomf.is" {
		return r.Files[0].URL
	}

	return h.ReturnURL + r.Files[0].URL
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
