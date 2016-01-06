package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/kori/gone/svc/teknik.io"
)

func main() {
	flag.Parse()
	var wg sync.WaitGroup

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: gone [files]")
		os.Exit(1)
	} else {
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
}
