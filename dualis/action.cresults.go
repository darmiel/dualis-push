package dualis

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"log"
	"strings"
)

type Course struct {
	NR         string
	Name       string
	EndGrade   string
	Credits    string
	Status     string
	ResultsURL string
}

func parseCourses(content string) (courses []*Course, err error) {
	body := strings.NewReader(content)

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(body); err != nil {
		return
	}

	doc.Find("tbody > tr").Each(func(i int, selection *goquery.Selection) {
		td := selection.Find("td")
		if len(td.Nodes) != 7 {
			return
		}

		values := make([]string, len(td.Nodes))
		td.Each(func(i int, selection *goquery.Selection) {
			t := selection.Text()
			if i == 5 { // link #5 contains the results
				// "pArSe" url for resuts
				// what the fuckkkkk
				x := strings.Index(t, "/scripts/")
				if x < 0 {
					return
				}
				t = t[x:]
				x = strings.Index(t, "\"")
				if x < 0 {
					return
				}
				t = t[:x]
			}
			values[i] = strings.TrimSpace(t)
		})

		courses = append(courses, &Course{
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

type Grade struct {
	ID       string
	Course   string
	Topic    string
	Semester string
	Type     string
	Grade    string
}

func (c *Client) CourseResults() (err error) {
	var resp *req.Response
	url := app(CourseResultsPrg, c.ArgumentsFromRefresh())
	if resp, err = c.c.R().Get(url); err != nil {
		panic(err)
	}

	var courses []*Course
	if courses, err = parseCourses(resp.String()); err != nil {
		return
	}

	log.Printf("Parsed courses: %+v", courses)
	return
}
