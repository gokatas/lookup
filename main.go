// Lookup looks up IP addresses concurrently. Hosts are supplied as STDIN stream.
// Adapted from https://youtu.be/woCg2zaIVzQ
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type lookup struct {
	host  string
	addrs []string
	err   error
}

func main() {
	hosts := make(chan string)
	lookups := make(chan lookup)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			hosts <- input.Text()
		}
		if err := input.Err(); err != nil {
			log.Fatal(err)
		}
		close(hosts)
		wg.Done()
	}()

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for host := range hosts {
				addrs, err := net.LookupHost(host)
				lookups <- lookup{host: host, addrs: addrs, err: err}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(lookups)
	}()

	for l := range lookups {
		fmt.Println(l.host, l.addrs)
	}
}
