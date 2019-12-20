package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"v2ray.com/core"
	"v2ray.com/core/infra/conf"
)

type VmessLink struct {
	Add  string `json:"add,omitempty"`
	Aid  string `json:"aid,omitempty"`
	Host string `json:"host,omitempty"`
	ID   string `json:"id,omitempty"`
	Net  string `json:"net,omitempty"`
	Path string `json:"path,omitempty"`
	Port string `json:"port,omitempty"`
	Ps   string `json:"ps,omitempty"`
	Tls  string `json:"tls,omitempty"`
	Type string `json:"type,omitempty"`
}

func (v VmessLink) String() string {
	return fmt.Sprintf("%s|%s|%s  (%s)", v.Net, v.Add, v.Port, v.Ps)
}

func (v VmessLink) GenOutbound() (*core.OutboundHandlerConfig, error) {

	out := &conf.OutboundDetourConfig{}
	out.Tag = "proxy"
	out.Protocol = "vmess"
	out.MuxSettings = &conf.MuxConfig{}

	p := conf.TransportProtocol(v.Net)
	s := &conf.StreamConfig{
		Network:  &p,
		Security: v.Tls,
	}

	switch v.Net {
	case "tcp":
		s.TCPSettings = &conf.TCPConfig{}
		if v.Type == "" || v.Type == "none" {
			s.TCPSettings.HeaderConfig = json.RawMessage([]byte(`{ "type": "none" }`))
		} else {
			s.TCPSettings.HeaderConfig = json.RawMessage([]byte(fmt.Sprintf(`
			{
				"type": "http",
				"request": {
					"path": ["%s"],
					"headers": {
						"Host": ["%s"]
					}
				}
			}
			`, v.Path, v.Host)))
		}
	case "kcp":
		s.KCPSettings = &conf.KCPConfig{}
		s.KCPSettings.HeaderConfig = json.RawMessage([]byte(fmt.Sprintf(`{ "type": "%s" }`, v.Type)))
	case "ws":
		s.WSSettings = &conf.WebSocketConfig{}
		s.WSSettings.Path = v.Path
		s.WSSettings.Headers = map[string]string{
			"Host": v.Host,
		}
	case "h2", "http":
		s.HTTPSettings = &conf.HTTPConfig{
			Path: v.Path,
		}
		if v.Host != "" {
			s.HTTPSettings.Host = &conf.StringList{v.Host}
		}
	}

	if v.Tls == "tls" {
		s.TLSSettings = &conf.TLSConfig{Insecure: true}
	}

	out.StreamSetting = s
	oset := json.RawMessage([]byte(fmt.Sprintf(`{
  "vnext": [
    {
      "address": "%s",
      "port": %s,
      "users": [
        {
          "id": "%s",
          "alterId": %s,
          "security": "auto"
        }
      ]
    }
  ]
}`, v.Add, v.Port, v.ID, v.Aid)))
	out.Settings = &oset
	return out.Build()
}

func parseVmess() (*VmessLink, error) {

	b, err := base64.StdEncoding.DecodeString(vmess[8:])
	if err != nil {
		return nil, err
	}

	v := &VmessLink{}
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}

	return v, nil
}
