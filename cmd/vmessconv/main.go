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
	showN := flag.Bool("n", false, "show v2rayN/NG format")
	showRK := flag.Bool("r", false, "show shadowrocket format")
	showQ := flag.Bool("q", false, "show Quantumult format")
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

	if *showN {
		fmt.Println("V2rayN:", lk.LinkStr("ng"))
	}
	if *showRK {
		fmt.Println("ShadowRocket:", lk.LinkStr("rk"))
	}
	if *showQ {
		fmt.Println("Quantumult:", lk.LinkStr("quan"))
	}
	if !*showN && !*showRK && !*showQ {
		fmt.Println("V2rayN:", lk.LinkStr("ng"))
		fmt.Println()
		fmt.Println("ShadowRocket:", lk.LinkStr("rk"))
		fmt.Println()
		fmt.Println("Quantumult:", lk.LinkStr("quan"))
	}
}
