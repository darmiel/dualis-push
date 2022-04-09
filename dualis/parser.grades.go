package dualis

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

const GradeNotSet = "noch nicht gesetzt"

type exam struct {
	ID       string
	Course   string
	Topic    string
	Semester string
	Type     string
	Grade    string
}

func parseExamGrades(content string) (ex []*exam, err error) {
	body := strings.NewReader(content)

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return
	}

	var e *exam

	doc.Find("tr").Each(func(i int, selection *goquery.Selection) {
		td := selection.Find("td")

		switch len(td.Nodes) {
		case 1:
			e = nil // reset header
			td.Each(func(i int, selection *goquery.Selection) {
				// exam headers have 1 td element with the "level02" class
				// why? only god knows...
				if !selection.HasClass("level02") {
					return
				}

				// the header is something like this:
				// ABCDEFGH.X ABC-DEFXX Analysis I
				spl := strings.SplitN(selection.Text(), " ", 3)
				if len(spl) == 3 {
					e = &exam{
						ID:     spl[0],
						Course: spl[1],
						Topic:  spl[2],
					}
				} else {
					e = &exam{
						ID:     "unknown",
						Course: "unknown",
						Topic:  spl[0],
					}
				}
			})

		case 6:
			// a exam header is required for exam data
			if e == nil {
				return
			}
			values := make([]string, 6)
			td.Each(func(i int, selection *goquery.Selection) {
				if !selection.HasClass("tbdata") {
					return
				}
				values[i] = strings.TrimSpace(selection.Text())
			})
			e.Semester = values[0]
			e.Type = values[1]
			// values[2] is always empty
			e.Grade = values[3]
			// values[4, 5] are always empty too...

			if e.Grade != GradeNotSet {
				// make a copy of e and append to ex
				ex = append(ex, &(*e))
			}
		}
	})

	return
}
