package dualis

import (
	"github.com/imroc/req/v3"
)

func (c *Client) CourseResults() (ex Grades, err error) {
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
		if resp, err = c.c.R().Get(semester.URL); err != nil {
			return
		}

		var modules []*semesterModule
		if modules, err = parseSemesterModules(resp.String()); err != nil {
			return
		}

		for _, module := range modules {
			if resp, err = c.c.R().Get(module.ResultsURL); err != nil {
				return
			}

			var exams []*Grade
			if exams, err = parseExamGrades(resp.String()); err != nil {
				return
			}

			ex = append(ex, exams...)
		}
	}

	return
}
