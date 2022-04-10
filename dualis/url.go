package dualis

import "fmt"

type URL string

const (
	DualisBaseURL   URL = "https://dualis.dhbw.de"
	DualisScriptURL     = DualisBaseURL + "/scripts/mgrqispi.dll"
)

const (
	ExternalPagesPrg = "EXTERNALPAGES"
	CourseResultsPrg = "COURSERESULTS"
)

func (u URL) WithArguments(args string) string {
	return string(u) + fmt.Sprintf("&ARGUMENTS=%s", args)
}

func App(name, args string) string {
	return fmt.Sprintf("%s?APPNAME=CampusNet&PRGNAME=%s&ARGUMENTS=%s", DualisScriptURL, name, args)
}
