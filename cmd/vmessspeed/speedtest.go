package main

import (
	"log"
	"os"
	"time"

	mv2ray "github.com/v2fly/vmessping/miniv2ray"
	"gopkg.in/alecthomas/kingpin.v2"
)

func checkError(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func setTimeout() {
	if *timeoutOpt != 0 {
		timeout = *timeoutOpt
	}
}

var (
	MAINVER = "0.0.0-src"
	timeout = 180

	vmessLink     = kingpin.Arg("vmess", "The link of VMess").Required().String()
	showList      = kingpin.Flag("list", "Show available speedtest.net servers").Short('l').Bool()
	debug         = kingpin.Flag("debug", "Show V2Ray core debug log").Short('d').Bool()
	serverIds     = kingpin.Flag("server", "Select server id to speedtest").Short('s').Ints()
	timeoutOpt    = kingpin.Flag("timeout", "Define timeout seconds. Default: 10 sec").Short('t').Int()
	useMux        = kingpin.Flag("mux", "Use Mux outbound").Short('m').Bool()
	allowInsecure = kingpin.Flag("allow-insecure", "Allow insecure TLS connections").Bool()
)

func main() {
	kingpin.Version(MAINVER)
	kingpin.Parse()

	setTimeout()

	server, err := mv2ray.StartV2Ray(*vmessLink, *debug, *useMux, *allowInsecure)
	if err != nil {
		log.Fatalln(err)
	}

	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
	defer server.Close()

	client, err = mv2ray.CoreHTTPClient(server, time.Second*time.Duration(timeout))
	if err != nil {
		log.Fatalln(err)
	}

	user := fetchUserInfo()
	user.Show()

	list := fetchServerList(user)
	if *showList {
		list.Show()
		return
	}

	targets := list.FindServer(*serverIds)
	targets.StartTest()
	targets.ShowResult()
}
