package dualis

import (
	"encoding/json"
	"github.com/apex/log"
)

type (
	Grades []*Grade
	Grade  struct {
		CourseID   string
		CourseName string
		Semester   string

		Grade string
		Type  string
	}
)

func (g *Grade) Marshal() (res string, err error) {
	var data []byte
	if data, err = json.Marshal(g); err != nil {
		return
	}
	res = string(data)
	return
}

func (g Grades) CheckForNewIn(other Grades, con func(g *Grade)) {
	for _, n := range other {
		// check if n in g
		nid, err := n.Marshal()
		if err != nil {
			log.Warn("cannot compare")
			return
		}

		found := false

	a:
		for _, o := range g {
			oid, err := o.Marshal()
			if err != nil {
				log.Warn("cannot compare")
				return
			}
			if oid == nid {
				found = true
				break a
			}
		}

		if !found {
			con(n) // send update
		}
	}
}
