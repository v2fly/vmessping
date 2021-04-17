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
	showS := flag.Bool("s", false, "show Standard format")
	showN := flag.Bool("n", false, "show v2rayN / v2rayNG format")
	showRK := flag.Bool("r", false, "show Shadowrocket format")
	showQ := flag.Bool("q", false, "show Quantumult format")
	flag.Parse()
	var link string
	if flag.NArg() == 0 {
		if link = os.Getenv("VMESS"); link == "" {
			fmt.Println(os.Args[0], "vmess:// ...")
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
	fmt.Println("VMessConvert:", MAINVER)

	if *showS {
		vlink, err := lk.LinkStr("s")
		if err != nil {
			fmt.Println("Standard:", err)
		}
		fmt.Println("Standard:", vlink)
	}

	if *showN {
		vlink, err := lk.LinkStr("ng")
		if err != nil {
			fmt.Println("v2rayN / v2rayNG:", err)
		}
		fmt.Println("v2rayN / v2rayNG:", vlink)
	}
	if *showRK {
		vlink, err := lk.LinkStr("rk")
		if err != nil {
			fmt.Println("Shadowrocket:", err)
		}
		fmt.Println("Shadowrocket:", vlink)
	}
	if *showQ {
		vlink, err := lk.LinkStr("quan")
		if err != nil {
			fmt.Println("Quantumult:", err)
		}
		fmt.Println("Quantumult:", vlink)
	}
	if !*showS && !*showN && !*showRK && !*showQ {
		var vlink string
		var err error

		vlink, err = lk.LinkStr("s")
		if err != nil {
			fmt.Println("Standard:", err)
		}
		fmt.Println("Standard:", vlink)

		vlink, err = lk.LinkStr("ng")
		if err != nil {
			fmt.Println("v2rayN / v2rayNG:", err)
		}
		fmt.Println("v2rayN / v2rayNG:", vlink)
		fmt.Println()

		vlink, err = lk.LinkStr("rk")
		if err != nil {
			fmt.Println("Shadowrocket:", err)
		}
		fmt.Println("Shadowrocket:", vlink)
		fmt.Println()

		vlink, err = lk.LinkStr("quan")
		if err != nil {
			fmt.Println("Quantumult:", err)
		}
		fmt.Println("Quantumult:", vlink)
	}
}
