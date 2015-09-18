package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/kori/gone/s/teknik"
)

var (
	www  = flag.String("w", "", "Website to upload to.")
	ctrs = flag.Bool("c", true, "Whether to display counters.")
)

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	c := 1 // current file counter

	// t = total files, p = path to file
	for t, p := range flag.Args() {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			f := teknik.Upload(p)
			if *ctrs {
				fmt.Print(c, "/", t+1, ": ")
			}
			fmt.Println(f)
			c++
		}(p)
	}
	wg.Wait()
}
