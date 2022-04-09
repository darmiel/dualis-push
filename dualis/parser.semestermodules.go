package dualis

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type semesterModule struct {
	NR         string
	Name       string
	EndGrade   string
	Credits    string
	Status     string
	ResultsURL string
}

func parseSemesterModules(content string) (sm []*semesterModule, err error) {
	body := strings.NewReader(content)

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return
	}

	// we just have to pray that there's only one table...
	doc.Find("tbody > tr").Each(func(i int, selection *goquery.Selection) {
		td := selection.Find("td")

		// oh wait! there ARE other tables.
		// let's just assume our important table is the only one with only 7 td nodes.
		if len(td.Nodes) != 7 {
			return
		}

		values := make([]string, len(td.Nodes))
		td.Each(func(i int, selection *goquery.Selection) {
			t := selection.Text()

			// the 6th td element contains a script that looks something like this:
			// ... = dl_popUp("/scripts/mgrqispi.dll?APPNAME=...&ARGUMENTS=$1,$2,$3","Resultdetails",750,400);
			// we need $1, $2 and $3
			if i == 5 {
				// /scripts/mgrqispi.dll?APPNAME=...&ARGUMENTS=$1,$2,$3","Resultdetails",750,400);
				x := strings.Index(t, "/scripts/")
				if x < 0 {
					return
				}
				t = t[x:]
				// /scripts/mgrqispi.dll?APPNAME=...&ARGUMENTS=$1,$2,$3
				x = strings.Index(t, "\"")
				if x < 0 {
					return
				}
				t = string(DualisBaseURL) + t[:x]
			}
			values[i] = strings.TrimSpace(t)
		})

		sm = append(sm, &semesterModule{
			NR:         values[0],
			Name:       values[1],
			EndGrade:   values[2],
			Credits:    values[3],
			Status:     values[4],
			ResultsURL: values[5],
		})
	})

	return
}
