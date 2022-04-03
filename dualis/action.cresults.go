package dualis

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"strings"
)

func (c *Client) CourseResults() (err error) {
	var resp *req.Response
	url := app(CourseResultsPrg, c.ArgumentsFromRefresh())
	if resp, err = c.c.R().Get(url); err != nil {
		panic(err)
	}

	var doc *goquery.Document
	body := strings.NewReader(resp.String())
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return
	}

	_ = doc // placeholder

	return
}
