package dualis

import (
	"errors"
	"github.com/imroc/req/v3"
	"strings"
)

var (
	ErrLoginNotSuccessful = errors.New("login not successful")
)

type Client struct {
	c       *req.Client
	refresh string
}

func (c *Client) ArgumentsFromRefresh() string {
	idx := strings.Index(c.refresh, "&ARGUMENTS=")
	if idx == -1 {
		return ""
	}
	str := c.refresh[idx+11:]
	if idx = strings.Index(str, "&"); idx != -1 {
		str = str[:idx]
	}
	return str[:strings.LastIndex(str, ",")+1]
}

// fixCookies fixes the retarded "Set-cookie" header sent by dualis
// PR: https://github.com/golang/go/pull/52121 should fix this,
// but it hasn't been accepted (yet)
func fixCookies(req *req.Response) {
	c := req.Header["Set-Cookie"]
	for i, v := range c {
		name, value, ok := strings.Cut(v, "=")
		if !ok {
			continue
		}
		if nt := strings.TrimSpace(name); name != nt {
			c[i] = nt + "=" + value
		}
	}
	req.Header["Set-Cookie"] = c
}
