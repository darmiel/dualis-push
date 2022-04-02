package main

import (
	"errors"
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"os"
)

type DualisURL string

const (
	DualisBaseURL   DualisURL = "https://dualis.dhbw.de"
	DualisScriptURL           = DualisBaseURL + "/scripts/mgrqispi.dll"
)

const (
	ExternalPagesPrg = "EXTERNALPAGES"
	CourseResultsPrg = "COURSERESULTS"
)

func (u DualisURL) WithArguments(args string) string {
	return string(u) + fmt.Sprintf("&ARGUMENTS=%s", args)
}

func app(name, args string) string {
	return fmt.Sprintf("%s?APPNAME=CampusNet&PRGNAME=%s&ARGUMENTS=%s", DualisScriptURL, name, args)
}

var (
	ErrLoginNotSuccessful = errors.New("login not successful")
)

type DualisClient struct {
	c       *req.Client
	refresh string
}

func main() {
	user, pass := os.Getenv("USER"), os.Getenv("PASS")
	fmt.Println("User:", user, "Pass:", pass)

	log.Println("Logging in...")
	client, err := Login(user, pass)
	if err != nil {
		panic(err)
	}
	log.Println("Logged in.")

	client.CourseResults()
}

func Login(username, password string) (client *DualisClient, err error) {
	c := req.C()
	var cookieReq *req.Response
	if cookieReq, err = c.R().Get(app(ExternalPagesPrg, "-N000000000000001,-N000324,-Awelcome")); err != nil {
		return
	}
	c.SetCommonCookies(cookieReq.Cookies()...)

	log.Printf("Headers: %+v", cookieReq.Header)

	// make login request
	loginPayload := map[string]string{
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
	var loginReq *req.Response
	if loginReq, err = c.R().
		SetFormData(loginPayload).
		Post(string(DualisScriptURL)); err != nil {
		return
	}

	refresh := loginReq.Header.Get("REFRESH")
	if refresh == "" {
		return nil, ErrLoginNotSuccessful
	}

	c.SetCommonCookies(loginReq.Cookies()...)
	client = &DualisClient{
		c:       c,
		refresh: refresh,
	}

	log.Printf("Login Res (%d): %+v (c), %+v (h)", loginReq.StatusCode, loginReq.Cookies(), loginReq.Header)
	return
}

func (c *DualisClient) CourseResults() {
	log.Println("Loading course results w/ refresh:", c.refresh, "...")
}
