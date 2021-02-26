package miniv2ray

import (
	"errors"
	"net/http"
	"strings"
	"time"

	core "github.com/v2fly/v2ray-core/v4"
)

const (
	cloudflareCGI = "http://www.cloudflare.com/cdn-cgi/trace"
)

func GetNodeInfo(inst *core.Instance, timeout time.Duration) (map[string]string, error) {
	code, bf, err := CoreHTTPRequest(inst, timeout, "GET", cloudflareCGI)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, errors.New("fail to get cdn-cgi/trace")
	}

	info := make(map[string]string)
	for _, line := range strings.Split(string(bf), "\n") {
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			continue
		}
		info[kv[0]] = kv[1]
	}

	if len(info) == 0 {
		return nil, errors.New("not getting anything from cdn-cgi/trace")
	}

	return info, nil
}
