module github.com/v2fly/vmessping

go 1.13

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/google/go-cmp v0.5.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	v2ray.com/core v4.19.1+incompatible
)

replace v2ray.com/core => github.com/v2fly/v2ray-core v1.24.5-0.20200531043819-9dc12961fac5
