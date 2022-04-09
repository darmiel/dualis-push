package dualis

import (
	"fmt"
	"github.com/imroc/req/v3"
	"log"
)

func (c *Client) CourseResults() (err error) {
	var resp *req.Response
	url := App(CourseResultsPrg, c.ArgumentsFromRefresh())
	if resp, err = c.c.R().Get(url); err != nil {
		return
	}

	// parse available semesters
	var semesters []*semesterOption
	if semesters, err = parseSemesterOptions(resp.String()); err != nil {
		return
	}

	for _, semester := range semesters {
		log.Println("Found Semester", semester.Name)

		if resp, err = c.c.R().Get(semester.URL); err != nil {
			return
		}

		var modules []*semesterModule
		if modules, err = parseSemesterModules(resp.String()); err != nil {
			return
		}

		for _, module := range modules {
			log.Println("  > Module:", module.Name)

			if resp, err = c.c.R().Get(module.ResultsURL); err != nil {
				return
			}

			var exams []*exam
			if exams, err = parseExamGrades(resp.String()); err != nil {
				return
			}

			for _, ex := range exams {
				log.Println("    > Topic:", ex.Topic)
				log.Println("    > Grade:", ex.Grade)
			}
		}

		fmt.Println()
	}

	return
}
