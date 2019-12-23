# vmessping
a ping prober for `vmess://` link in V2rayN/NG format.

```
./vmessping vmess://....
Usage of ./vmessping:
  -c uint
        Count. Stop after sending COUNT requests (default 9999)
  -dest string
        the test destination url, need 204 for success return (default "http://www.google.com/gen_204")
  -i uint
        inteval seconds between pings (default 1)
  -o uint
        timeout seconds for each request (default 10)
  -v    verbose (debug log)
```

# Example
```
./vmessping -c 3 vmess://12345678.......

Vmessping [0.0.0-src] Yet another distribution of v2ray (v2ray-core: 4.21.3)
PING  tcp|my.server.domian|4321  (ps/name)

2019/12/20 15:56:09 Get http://www.google.com/gen_204: net/http: request canceled ...
Ping http://www.google.com/gen_204: seq=1 time=-1 ms
Ping http://www.google.com/gen_204: seq=2 time=490 ms
Ping http://www.google.com/gen_204: seq=3 time=396 ms

--- http://www.google.com/gen_204 vmess ping statistics ---
3 requests made, 2 success, time 6.886747334s
rtt min/avg/max = 396/443.00/490 ms
```

# Compile from source
```
go get -v github.com/v2fly/vmessping/vmessping/...
```