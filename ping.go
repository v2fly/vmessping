package vmessping

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"v2ray.com/core"
)

func PrintVersion(mv string) {
	fmt.Fprintf(os.Stderr,
		"Vmessping [%s] Yet another distribution of v2ray (v2ray-core: %s)\n", mv, core.Version())
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

func Ping(vmess string, count uint, dest string, timeoutsec, inteval uint, verbose bool) (int, error) {
	server, err := StartV2Ray(vmess, verbose)
	if err != nil {
		fmt.Println(err.Error())
		return 2, err
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		return 2, err
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
		delay, err := MeasureDelay(server, time.Second*time.Duration(timeoutsec), dest)
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
			time.Sleep(time.Second * time.Duration(inteval))
		}
	}

	printStat(delays, reqcnt, errcnt)
	return 0, nil
}
