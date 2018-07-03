package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/kori/gone/services/teknik"
)

func main() {
	var wg sync.WaitGroup

	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatalln("Usage: gone [files]")
	}

	// Upload each file.
	for _, p := range flag.Args() {
		wg.Add(1)
		go func(file string) {
			if exists(file) {
				link, err := teknik.Upload(file)
				if err != nil {
					fmt.Println(file+":", "upload failed:", err)
				}
				fmt.Println(file+":", link)
			} else {
				fmt.Println("file doesn't exist:", file)
			}
			wg.Done()
		}(p)
	}

	wg.Wait()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
