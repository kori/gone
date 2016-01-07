package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kori/gone/svc/pomf"
	"github.com/kori/gone/svc/teknik.io"
)

var (
	host    = flag.String("h", "teknik", "Which type of host to use. Current options: teknik/pomf")
	service = flag.String("s", "", "Which service to use. Current options, teknik: upload, pomf: 1339.cf, maxfile.ro, pomf.cat, mixtape.moe, pomf.is")
)

func main() {
	flag.Parse()
	var wg sync.WaitGroup

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: gone -h [host] -s [service] [files]")
		os.Exit(1)
	} else {
		// Upload each file.
		for _, p := range flag.Args() {
			wg.Add(1)
			go func(file string) {
				defer wg.Done()
				if exists(file) {
					switch *host {
					case "pomf":
						fmt.Println(pomf.Upload(file, *service))
					case "teknik":
						fmt.Println(teknik.Upload(file))
					default:
						log.Fatal("no host provided")
					}
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
