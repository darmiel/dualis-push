package dualis

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

const GradeNotSet = "noch nicht gesetzt"

var (
	spaceRe = regexp.MustCompile(`\s+`)
	semRe   = regexp.MustCompile(`\((.*)\)`)
)

func parseExamGrades(content string) (ex Grades, err error) {
	body := strings.NewReader(content)

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return
	}

	// find course ID and course name
	var (
		courseID   string
		courseName string
		semester   string
	)

	// warning: cursed
	{
		h1 := strings.TrimSpace(doc.Find("h1").First().Text())
		x := strings.Index(h1, "(")
		if x > 0 {
			semester = strings.Trim(semRe.FindString(h1), "()")
			h1 = h1[:x]
		}
		spl := spaceRe.Split(h1, 2)
		courseID = strings.TrimSpace(spl[0])
		courseName = strings.TrimSpace(spl[1])
	}

	ex = make(Grades)
	doc.Find("tr").Each(func(i int, selection *goquery.Selection) {
		td := selection.Find("td")

		if len(td.Nodes) == 6 {
			var values []string
			td.Each(func(i int, selection *goquery.Selection) {
				if !selection.HasClass("tbdata") {
					return
				}
				values = append(values, strings.TrimSpace(selection.Text()))
			})
			if len(values) != 6 {
				return
			}
			grade := values[3]
			if grade == "" || grade == GradeNotSet {
				return
			}

			g := &Grade{
				CourseID:   courseID,
				CourseName: courseName,
				Semester:   semester,
				Grade:      grade,
				Type:       values[1],
			}
			ex[g.Marshal()] = g
		}
	})

	return
}
