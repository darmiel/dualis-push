package dualis

import (
	"encoding/json"
	"strconv"
	"strings"
)

type (
	Grades map[string]*Grade
	Grade  struct {
		CourseID   string
		CourseName string
		Semester   string

		Grade string
		Type  string
	}
)

func (g *Grade) Unit() string {
	graw := g.Grade
	if strings.Contains(graw, ",") {
		graw = graw[:strings.Index(graw, ",")]
	}
	if grawi, err := strconv.Atoi(graw); err == nil {
		if grawi <= 5 && grawi >= 1 {
			return "Note"
		} else {
			return "%"
		}
	}
	return "Pommes"
}

func (g *Grade) Marshal() string {
	data, _ := json.Marshal(g)
	return string(data)
}
