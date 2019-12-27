package main

import (
	"flag"
	"fmt"
	"github.com/v2fly/vmessping"
	"os"
)

var (
	MAINVER = "0.0.0-src"
)

func main() {
	verbose := flag.Bool("v", false, "verbose (debug log)")
	desturl := flag.String("dest", "http://www.google.com/gen_204", "the test destination url, need 204 for success return")
	count := flag.Uint("c", 9999, "Count. Stop after sending COUNT requests")
	timeout := flag.Uint("o", 10, "timeout seconds for each request")
	inteval := flag.Uint("i", 1, "inteval seconds between pings")
	quit := flag.Uint("q", 0, "fast quit on error counts")
	flag.Parse()

	var vmess string
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
	code, err := vmessping.Ping(vmess, *count, *desturl, *timeout, *inteval, *quit, *verbose)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}
