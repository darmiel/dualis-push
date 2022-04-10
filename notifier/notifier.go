package notifier

import (
	"errors"
	"github.com/darmiel/dualis-push/dualis"
	"strings"
)

type (
	Func      func(notifier *Notifier, grade *dualis.Grade, repl Replacers) error
	Replacers map[string]string
)

var (
	ErrUnknownNotifierType = errors.New("unknown notifier")
	notifiers              = make(map[string]Func)
)

type Notifier struct {
	Type     string
	Disabled bool

	// Pushover
	PushoverToken     string
	PushoverRecipient string

	// Discord
	DiscordWebhookURL string

	Format *Format
}

func replace(message string, replacers Replacers) string {
	for k, v := range replacers {
		message = strings.Replace(message, "%"+strings.ToLower(k)+"%", v, -1)
	}
	return message
}

func (n *Notifier) Formatting() Format {
	if n.Format != nil {
		return *n.Format
	}
	return DefaultFormat
}

func (n *Notifier) DoGradeArrived(grade *dualis.Grade) error {
	f, ok := notifiers[strings.ToLower(n.Type)]
	if !ok {
		return ErrUnknownNotifierType
	}
	return f(n, grade, Replacers{
		"grade":     grade.Grade,
		"course":    grade.CourseName,
		"course-id": grade.CourseID,
		"semester":  grade.Semester,
	})
}
