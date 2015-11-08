package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/kori/gone/svc/teknik"
)

func main() {
	flag.Parse()
	var wg sync.WaitGroup

	// p = path to file
	for _, p := range flag.Args() {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			if _, err := os.Stat(file); err == nil {
				url := teknik.Upload(p)
				fmt.Println(url)
			} else {
				fmt.Println(file, "doesn't exist")
			}
		}(p)
	}
	wg.Wait()
}
