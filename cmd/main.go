package main

import (
	"flag"
	"fmt"
	"github.com/v2fly/vmessping"
	"os"
)

var (
	vmess   string
	count   uint
	timeout uint
	desturl string
	MAINVER = "0.0.0-src"
)

func main() {
	verbose := flag.Bool("v", false, "verbose (debug log)")
	flag.StringVar(&desturl, "dest", "http://www.google.com/gen_204", "the test destination url, need 204 for success return")
	flag.UintVar(&count, "c", 9999, "Count. Stop after sending COUNT requests")
	flag.UintVar(&timeout, "o", 10, "timeout seconds for each request")
	flag.Parse()

	if flag.NArg() == 0 {
		if vmess = os.Getenv("VMESS"); vmess == "" {
			fmt.Println(os.Args[0], "vmess://....")
			flag.Usage()
			os.Exit(1)
		}
	} else {
		vmess = flag.Args()[0]
	}

	vmessping.PrintVersion(MAINVER)
	code, err := vmessping.Ping(vmess, count, desturl, *verbose)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}
