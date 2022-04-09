package dualis

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type semesterOption struct {
	Name string
	URL  string
}

func parseSemesterOptions(content string) (so []*semesterOption, err error) {
	body := strings.NewReader(content)

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return
	}

	// the selector contains available semesters
	sel := doc.Find("select").First()

	// parse fucking URL
	// there is absolutely no way this API is going to be compatible for some time
	// really hacky as I don't really want to invest any longer in parsing this crap.

	// anyway, the onchange attribute contains an event that looks something like this:
	// `reloadpage.createUrlAndReload('/scripts/mgrqispi.dll','CampusNet','COURSERESULTS','$1','$2','-N'+this.value);`
	url, _ := sel.Attr("onchange")

	// we only want go get $1 and $2
	// splitting at ',' $1 should be at index 3, $2 at index 4
	us := strings.Split(url, ",")

	// create base URL which looks something like this:
	// https://dualis.dhbw.de/scripts/mgrqispi.dll?APPNAME=CampusNet&PRGNAME=COURSERESULTS&ARGUMENTS=-N$1,-N$2,-N$3
	baseURL := App("COURSERESULTS", strings.Join([]string{
		"-N" + us[3][1:len(us[3])-1], "-N" + us[4][1:len(us[4])-1],
	}, ","))

	// "parse" semesters with corresponding URLs
	sel.Find("option").Each(func(i int, selection *goquery.Selection) {
		param, _ := selection.Attr("value")
		so = append(so, &semesterOption{
			Name: selection.Text(),
			URL:  baseURL + ",-N" + param,
		})
	})

	return
}
