package vmessping

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"v2ray.com/core"
	v2net "v2ray.com/core/common/net"
)

func PrintVersion(mv string) {
	fmt.Fprintf(os.Stderr,
		"Vmessping [%s] Yet another distribution of v2ray (v2ray-core: %s)\n", mv, core.Version())
}

func MeasureDelay(inst *core.Instance, timeout time.Duration, dest string) (int64, error) {
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
		Timeout:   timeout,
	}

	req, _ := http.NewRequest("GET", dest, nil)
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

func printStat(delays []int64, req, errs int) {
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

	fmt.Println("\n--- vmess ping statistics ---")
	fmt.Printf("%d requests made, %d success, total time %v\n", req, len(delays), time.Duration(sum)*time.Millisecond)
	fmt.Printf("rtt min/avg/max = %d/%.0f/%d ms\n", min, avg, max)
}

func Ping(vmess string, count uint, dest string, verbose bool) (int, error) {
	server, err := StartV2Ray(vmess, verbose)
	if err != nil {
		fmt.Println(err.Error())
		return -1, err
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		return -1, err
	}
	defer server.Close()

	round := count
	var delays []int64
	var errcnt int
	var reqcnt int

	go func() {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-osSignals
		fmt.Println()
		printStat(delays, reqcnt, errcnt)
		if len(delays) == 0 {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	for round > 0 {
		seq := count - round + 1

		reqcnt++
		delay, err := MeasureDelay(server, time.Second*10, dest)
		if err != nil {
			errcnt++
		}

		if delay > 0 {
			delays = append(delays, delay)
			fmt.Printf("Ping %s: seq=%d time=%d ms\n", dest, seq, delay)
		} else {
			fmt.Printf("Ping %s: seq=%d err %v\n", dest, seq, err)
		}

		round--
		if round > 0 {
			time.Sleep(time.Second)
		}
	}

	printStat(delays, reqcnt, errcnt)
	return 0, nil
}
