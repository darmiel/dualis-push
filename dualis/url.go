package dualis

import "fmt"

type DualisURL string

const (
	DualisBaseURL   DualisURL = "https://dualis.dhbw.de"
	DualisScriptURL           = DualisBaseURL + "/scripts/mgrqispi.dll"
)

const (
	ExternalPagesPrg = "EXTERNALPAGES"
	CourseResultsPrg = "COURSERESULTS"
)

func (u DualisURL) WithArguments(args string) string {
	return string(u) + fmt.Sprintf("&ARGUMENTS=%s", args)
}

func App(name, args string) string {
	return fmt.Sprintf("%s?APPNAME=CampusNet&PRGNAME=%s&ARGUMENTS=%s", DualisScriptURL, name, args)
}
