package vmessping

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	applog "v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	commlog "v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/infra/conf"
)

type VmessLink struct {
	Add  string      `json:"add,omitempty"`
	Aid  string      `json:"aid,omitempty"`
	Host string      `json:"host,omitempty"`
	ID   string      `json:"id,omitempty"`
	Net  string      `json:"net,omitempty"`
	Path string      `json:"path,omitempty"`
	Port interface{} `json:"port,omitempty"`
	Ps   string      `json:"ps,omitempty"`
	TLS  string      `json:"tls,omitempty"`
	Type string      `json:"type,omitempty"`
}

func (v VmessLink) String() string {
	return fmt.Sprintf("%s|%s|%v  (%s)", v.Net, v.Add, v.Port, v.Ps)
}

func (v VmessLink) GenOutbound() (*core.OutboundHandlerConfig, error) {

	out := &conf.OutboundDetourConfig{}
	out.Tag = "proxy"
	out.Protocol = "vmess"
	out.MuxSettings = &conf.MuxConfig{}

	p := conf.TransportProtocol(v.Net)
	s := &conf.StreamConfig{
		Network:  &p,
		Security: v.TLS,
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

	if v.TLS == "tls" {
		s.TLSSettings = &conf.TLSConfig{Insecure: true}
	}

	out.StreamSetting = s
	oset := json.RawMessage([]byte(fmt.Sprintf(`{
  "vnext": [
    {
      "address": "%s",
      "port": %v,
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

func NewVmess(vmess string) (*VmessLink, error) {

	b64 := vmess[8:]
	if len(b64)/4 != 0 {
		b64 += strings.Repeat("=", len(b64)%4)
	}

	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	v := &VmessLink{}
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}

	return v, nil
}

func StartV2Ray(vmess string, verbose bool) (*core.Instance, error) {

	loglevel := commlog.Severity_Error
	if verbose {
		loglevel = commlog.Severity_Debug
	}

	o, err := NewVmess(vmess)
	if err != nil {
		return nil, err
	}

	fmt.Println("PING ", o.String())
	ob, err := o.GenOutbound()
	if err != nil {
		return nil, err
	}
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&applog.Config{
				ErrorLogType:  applog.LogType_Console,
				ErrorLogLevel: loglevel,
			}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	}

	commlog.RegisterHandler(commlog.NewLogger(commlog.CreateStderrLogWriter()))
	config.Outbound = []*core.OutboundHandlerConfig{ob}
	server, err := core.New(config)
	if err != nil {
		return nil, err
	}

	return server, nil
}
