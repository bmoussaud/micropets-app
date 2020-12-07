package main

import (
	"fmt"
	"net"
	"os"
)

func index(service string) {

	fmt.Fprintf(os.Stderr, "Service %v\n", service)
	ips, err := net.LookupIP(service)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get IPs: %v\n", err)
		os.Exit(1)
	}
	for _, ip := range ips {
		fmt.Printf("%s. IN A %s\n", service, ip.String())
	}

}

func main() {
	service := os.Args[1]

	index(service)
	index(service)
	index(service)
	index(service)
	index(service)
	index(service)

}
