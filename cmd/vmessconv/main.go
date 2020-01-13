package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/v2fly/vmessping/vmess"
)

var (
	MAINVER = "0.0.0-src"
)

func main() {
	flag.Parse()
	var link string
	if flag.NArg() == 0 {
		if link = os.Getenv("VMESS"); link == "" {
			fmt.Println(os.Args[0], "vmess://....")
			flag.Usage()
			os.Exit(1)
		}
	} else {
		link = flag.Args()[0]
	}

	lk, err := vmess.ParseVmess(link)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("VmessConvert:", MAINVER)
	fmt.Println("V2rayN:", lk.LinkStr("ng"))
	fmt.Println("ShadowRocket:", lk.LinkStr("rk"))
	fmt.Println("Quantumult:", lk.LinkStr("quan"))
}
