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

	// t = total files, p = path to file
	for t, p := range flag.Args() {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			url := teknik.Upload(p)
			fmt.Println(url)
		}(p)
	}
	wg.Wait()
}
