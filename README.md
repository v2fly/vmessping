# vmessping
![Build Status](https://github.com/v2fly/vmessping/workflows/Go/badge.svg) 

A ping prober for `vmess://` links in common seen formats.

`vmessping` supports the following fomats:

* V2rayN (V2rayNG)
* Shadowrocket 
* Quantumult (X)

# Download

Binaries are built automaticly by Github Action.

Download in [Release](https://github.com/v2fly/vmessping/releases/latest)

# Usage

```
./vmessping vmess://....
Usage of ./vmessping:
  -c uint
        Count. Stop after sending COUNT requests (default 9999)
  -dest string
        the test destination url, need 204 for success return (default "http://www.google.com/gen_204")
  -i uint
        inteval seconds between pings (default 1)
  -m    use mux outbound
  -n    show node location/outbound ip
  -o uint
        timeout seconds for each request (default 10)
  -q uint
        fast quit on error counts
  -v    verbose (debug log)
```

# Example
```
./vmessping "vmess://ew0KI......."
Vmessping ver[0.0.0-src], A prober for v2ray (v2ray-core: 4.23.1)

Type: ws
Addr: v2-server.address
Port: 80
UUID: 00000000-0000-0000-0000-000000000000
PS: @describe

Ping http://www.google.com/gen_204: seq=1 time=770 ms
Ping http://www.google.com/gen_204: seq=2 time=1368 ms
Ping http://www.google.com/gen_204: seq=3 time=761 ms
Ping http://www.google.com/gen_204: seq=4 time=761 ms
Ping http://www.google.com/gen_204: seq=5 time=1352 ms
^C
--- vmess ping statistics ---
5 requests made, 5 success, total time 9.106869693s
rtt min/avg/max = 761/1002/1368 ms
```

# Compile from source
```
git clone https://github.com/v2fly/vmessping.git
cd vmessping/cmd/vmessping
go build
```

# Other tools

## VmessSpeed

Speedtest for vmess.

### Usage
```

./vmessspeed --help
usage: vmessspeed [<flags>] <vmess>

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -l, --list               Show available speedtest.net servers
  -d, --debug              Show v2ray core debug log
  -s, --server=SERVER ...  Select server id to speedtest
  -t, --timeout=TIMEOUT    Define timeout seconds. Default: 10 sec
      --version            Show application version.

Args:
  <vmess>  the vmesslink

```

### Example

```
D:\>vmessspeed_amd64_windows.exe -s 14791 vmess://.....

Type: ws
Addr: v2-server.address
Port: 80
UUID: 00000000-0000-0000-0000-000000000000
PS: @describe

Testing From IP: ...IP..ADDR....

Target Server: [14791]    63.21km Macau (Macau) by MTel
Latency: 21.005ms
Download Test: ................11.97 Mbit/s
Upload Test: ................15.11 Mbit/s

Download: 11.97 Mbit/s
Upload: 15.11 Mbit/s
```