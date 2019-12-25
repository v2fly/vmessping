package vmess

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type VmessLink struct {
	Ver      string      `json:"v"`
	Add      string      `json:"add"`
	Aid      string      `json:"aid"`
	Host     string      `json:"host"`
	ID       string      `json:"id"`
	Net      string      `json:"net"`
	Path     string      `json:"path"`
	Port     interface{} `json:"port"`
	Ps       string      `json:"ps"`
	TLS      string      `json:"tls"`
	Type     string      `json:"type"`
	OrigLink string      `json:"-"`
}

func (v *VmessLink) IsEqual(c *VmessLink) bool {
	return v.Add == c.Add && v.Aid == c.Aid &&
		v.Host == c.Host && v.ID == c.ID &&
		v.Net == c.Net && v.Path == c.Path &&
		v.Port == c.Port && v.TLS == c.TLS &&
		v.Type == c.Type
}

func (v VmessLink) LinkStr() string {
	b, _ := json.Marshal(v)
	return "vmess://" + base64.URLEncoding.EncodeToString(b)
}

func (v VmessLink) String() string {
	return fmt.Sprintf("%s|%s|%v  (%s)", v.Net, v.Add, v.Port, v.Ps)
}

func NewVnVmess(vmess string) (*VmessLink, error) {

	if !strings.HasPrefix(vmess, "vmess://") {
		return nil, fmt.Errorf("vmess unreconized: %s", vmess)
	}

	b64 := vmess[8:]
	b, err := Base64Decode(b64)
	if err != nil {
		return nil, err
	}

	v := &VmessLink{}
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	v.OrigLink = vmess

	return v, nil
}

func NewRkVmess(vmess string) (*VmessLink, error) {
	if !strings.HasPrefix(vmess, "vmess://") {
		return nil, fmt.Errorf("vmess unreconized: %s", vmess)
	}
	url, err := url.Parse(vmess)
	if err != nil {
		return nil, err
	}
	link := &VmessLink{}

	b64 := url.Host
	b, err := Base64Decode(b64)
	if err != nil {
		return nil, err
	}

	mhp := strings.SplitN(string(b), ":", 3)
	if len(mhp) != 3 {
		return nil, fmt.Errorf("vmess unreconized: method:host:port -- %v", mhp)
	}
	link.Type = mhp[0]
	link.Port = mhp[2]
	idadd := strings.SplitN(mhp[1], "@", 2)
	if len(idadd) != 2 {
		return nil, fmt.Errorf("vmess unreconized: id@addr -- %v", idadd)
	}
	link.ID = idadd[0]
	link.Add = idadd[1]
	link.Aid = "0"

	vals := url.Query()
	if v := vals.Get("remarks"); v != "" {
		link.Ps = v
	}
	if v := vals.Get("path"); v != "" {
		link.Path = v
	}
	if v := vals.Get("tls"); v == "1" {
		link.TLS = "tls"
	}
	if v := vals.Get("obfs"); v != "" {
		switch v {
		case "websocket":
			link.Net = "ws"
		}
	}
	if v := vals.Get("obfsParam"); v != "" {
		link.Host = v
	}

	link.OrigLink = vmess
	return link, nil
}

func Base64Decode(b64 string) ([]byte, error) {
	if pad := len(b64) % 4; pad != 0 {
		b64 += strings.Repeat("=", 4-pad)
	}

	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return base64.URLEncoding.DecodeString(b64)
	}
	return b, nil
}

func ParseVmess(vmess string) (*VmessLink, error) {
	var lk *VmessLink
	if o, nerr := NewVnVmess(vmess); nerr == nil {
		lk = o
	} else if o, rerr := NewRkVmess(vmess); rerr == nil {
		lk = o
	} else {
		return nil, fmt.Errorf("%v, %v", nerr, rerr)
	}
	return lk, nil
}
