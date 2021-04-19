package vmessping

import (
	"fmt"
	"os"
	"time"

	mv2ray "github.com/v2fly/vmessping/miniv2ray"
)

func PrintVersion(mv string) {
	fmt.Fprintf(os.Stderr,
		"VMessPing ver[%s], A prober for v2ray (v2ray-core: %s)\n", mv, mv2ray.CoreVersion())
}

type PingStat struct {
	StartTime  time.Time
	SumMs      uint
	MaxMs      uint
	MinMs      uint
	AvgMs      uint
	Delays     []int64
	ReqCounter uint
	ErrCounter uint
}

func (p *PingStat) CalStats() {
	for _, v := range p.Delays {
		p.SumMs += uint(v)
		if p.MaxMs == 0 || p.MinMs == 0 {
			p.MaxMs = uint(v)
			p.MinMs = uint(v)
		}
		if uv := uint(v); uv > p.MaxMs {
			p.MaxMs = uv
		}
		if uv := uint(v); uv < p.MinMs {
			p.MinMs = uv
		}
	}
	if len(p.Delays) > 0 {
		p.AvgMs = uint(float64(p.SumMs) / float64(len(p.Delays)))
	}
}

func (p PingStat) PrintStats() {
	fmt.Println("\n--- vmess ping statistics ---")
	fmt.Printf("%d requests made, %d success, total time %v\n", p.ReqCounter, len(p.Delays), time.Since(p.StartTime))
	fmt.Printf("rtt min/avg/max = %d/%d/%d ms\n", p.MinMs, p.AvgMs, p.MaxMs)
}

func (p PingStat) IsErr() bool {
	return len(p.Delays) == 0
}

func Ping(vmess string, count uint, dest string, timeoutsec, inteval, quit uint, stopCh <-chan os.Signal, showNode, verbose, useMux, allowInsecure bool) (*PingStat, error) {
	server, err := mv2ray.StartV2Ray(vmess, verbose, useMux, allowInsecure)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		return nil, err
	}
	defer server.Close()

	if showNode {
		go func() {
			info, err := mv2ray.GetNodeInfo(server, time.Second*10)
			if err != nil {
				return
			}

			fmt.Printf("Node Outbound: %s/%s\n", info["loc"], info["ip"])
		}()
	}

	ps := &PingStat{}
	ps.StartTime = time.Now()
	round := count
L:
	for round > 0 {
		seq := count - round + 1
		ps.ReqCounter++

		chDelay := make(chan int64)
		go func() {
			delay, err := mv2ray.MeasureDelay(server, time.Second*time.Duration(timeoutsec), dest)
			if err != nil {
				ps.ErrCounter++
				fmt.Printf("Ping %s: seq=%d err %v\n", dest, seq, err)
			}
			chDelay <- delay
		}()

		select {
		case delay := <-chDelay:
			if delay > 0 {
				ps.Delays = append(ps.Delays, delay)
				fmt.Printf("Ping %s: seq=%d time=%d ms\n", dest, seq, delay)
			}
		case <-stopCh:
			break L
		}

		if quit > 0 && ps.ErrCounter >= quit {
			break
		}

		if round--; round > 0 {
			select {
			case <-time.After(time.Second * time.Duration(inteval)):
				continue
			case <-stopCh:
				break L
			}
		}
	}

	ps.CalStats()
	return ps, nil
}
