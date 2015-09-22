package main

import (
	"flag"
	"fmt"
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
			url := teknik.Upload(p)
			fmt.Println(url)
		}(p)
	}
	wg.Wait()
}
