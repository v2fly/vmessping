//go:generate go-bindata -nomemcopy -o bindata.go dat/...

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"v2ray.com/core"
	commlog "v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	_ "v2ray.com/core/main/distro/all"
)

var (
	vmess    string
	loglevel commlog.Severity
	multiple uint
	timeout  uint
	desturl  string = "http://www.google.com/gen_204"
	MAINVER         = "0.0.0-src"
)

func printVersion() {
	fmt.Fprintf(os.Stderr, "Vmessping [%s] Yet another distribution of v2ray (v2ray-core: %s)\n", MAINVER, core.Version())
}

func measureInstDelay(ctx context.Context, inst *core.Instance) (int64, error) {
	if inst == nil {
		return -1, newError("core instance nil")
	}

	tr := &http.Transport{
		TLSHandshakeTimeout: 6 * time.Second,
		DisableKeepAlives:   true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dest, err := v2net.ParseDestination(fmt.Sprintf("%s:%s", network, addr))
			if err != nil {
				return nil, err
			}
			return core.Dial(ctx, inst, dest)
		},
	}

	c := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", desturl, nil)
	start := time.Now()
	resp, err := c.Do(req)
	if err != nil {
		return -1, err
	}
	if resp.StatusCode != http.StatusNoContent {
		return -1, fmt.Errorf("status != 204: %s", resp.Status)
	}
	resp.Body.Close()
	return time.Since(start).Milliseconds(), nil
}

func printStat(delays []int64, req, errs int, start time.Time) {
	var sum int64
	var max int64
	var min int64
	for _, v := range delays {
		sum += v
		if max == 0 || min == 0 {
			max = v
			min = v
		}
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	avg := float64(sum) / float64(len(delays))

	fmt.Printf("\n--- %s vmess ping statistics ---\n", desturl)
	fmt.Printf("%d requests made, %d success, total time %v\n", req, len(delays), time.Since(start))
	fmt.Printf("rtt min/avg/max = %d/%.0f/%d ms\n", min, avg, max)
}

func main() {
	version := flag.Bool("version", false, "Show current version.")
	verbose := flag.Bool("v", false, "verbose (debug log)")
	flag.StringVar(&desturl, "dest", "http://www.google.com/gen_204", "the test destination url, need 204 for success return")
	flag.UintVar(&multiple, "c", 9999, "Count. Stop after sending COUNT requests")
	flag.UintVar(&timeout, "o", 10, "timeout seconds for each request")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println(os.Args[0], "vmess://....")
		flag.Usage()
		os.Exit(1)
	}

	printVersion()
	if *version {
		return
	}
	loglevel = commlog.Severity_Error
	if *verbose {
		loglevel = commlog.Severity_Debug
	}

	vmess = flag.Args()[0]
	server, err := startV2Ray()
	if err != nil {
		fmt.Println(err.Error())
		// Configuration error. Exit with a special value to prevent systemd from restarting.
		os.Exit(23)
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		os.Exit(-1)
	}
	defer server.Close()

	round := multiple
	var delays []int64
	var errcnt int
	var reqcnt int
	startTime := time.Now()

	go func() {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
		fmt.Println()
		printStat(delays, reqcnt, errcnt, startTime)
		if len(delays) == 0 {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	for round > 0 {
		reqcnt++
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		delay, err := measureInstDelay(ctx, server)
		cancel()
		if err != nil {
			errcnt++
			log.Println(err)
		}
		if delay > 0 {
			delays = append(delays, delay)
		}
		fmt.Printf("Ping %s: seq=%d time=%d ms\n", desturl, multiple-round+1, delay)
		round--
		time.Sleep(time.Second)
	}
	printStat(delays, reqcnt, errcnt, startTime)
	if len(delays) == 0 {
		os.Exit(1)
	}
}
