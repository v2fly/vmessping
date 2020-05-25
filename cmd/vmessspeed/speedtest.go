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
	vmessLink  = kingpin.Arg("vmess", "the vmesslink").Required().String()
	showList   = kingpin.Flag("list", "Show available speedtest.net servers").Short('l').Bool()
	debug      = kingpin.Flag("debug", "Show v2ray core debug log").Short('d').Bool()
	serverIds  = kingpin.Flag("server", "Select server id to speedtest").Short('s').Ints()
	timeoutOpt = kingpin.Flag("timeout", "Define timeout seconds. Default: 10 sec").Short('t').Int()
	timeout    = 180
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	setTimeout()

	server, err := mv2ray.StartV2Ray(*vmessLink, *debug, true)
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
