package dualis

import (
	"github.com/imroc/req/v3"
)

func Login(username, password string) (client *Client, err error) {
	c := req.C()

	// fix cookies
	c.OnAfterResponse(func(_ *req.Client, response *req.Response) error {
		fixCookies(response)
		return nil
	})

	// make login request
	var loginReq *req.Response
	if loginReq, err = c.R().
		SetFormData(loginPayload(username, password)).
		Post(string(DualisScriptURL)); err != nil {
		return
	}

	refresh := loginReq.Header.Get("REFRESH")
	if refresh == "" {
		return nil, ErrLoginNotSuccessful
	}

	c.SetCommonCookies(loginReq.Cookies()...)
	client = &Client{
		c:       c,
		refresh: refresh,
	}
	return
}

func loginPayload(username, password string) map[string]string {
	return map[string]string{
		"usrname":   username,
		"pass":      password,
		"APPNAME":   "CampusNet",
		"PRGNAME":   "LOGINCHECK",
		"ARGUMENTS": "clino,usrname,pass,menuno,menu_type, browser,platform",
		"clino":     "000000000000001",
		"menuno":    "000324",
		"menu_type": "classic",
		"browser":   "",
		"platform":  "",
	}
}
