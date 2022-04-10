package notifier

import (
	"errors"
	"github.com/darmiel/dualis-push/dualis"
	"github.com/imroc/req/v3"
	"strings"
)

var (
	ErrDiscordWebhookMissing = errors.New("webhook url missing")
)

func init() {
	type payload struct {
		Content string `json:"content"`
	}
	notifiers["discord"] = func(notifier *Notifier, grade *dualis.Grade, repl Replacers) error {
		if notifier.DiscordWebhookURL == "" {
			return ErrDiscordWebhookMissing
		}

		var (
			f   = notifier.Formatting()
			bob strings.Builder
		)
		if f.NewGradeMessageTitle != "" {
			bob.WriteString("**")
			bob.WriteString(f.NewGradeMessageTitle)
			bob.WriteString("**")
		}
		if f.NewGradeMessageBody != "" {
			if f.NewGradeMessageTitle != "" {
				bob.WriteRune('\n')
			}
			bob.WriteString(f.NewGradeMessageBody)
		}

		_, err := req.R().SetBodyJsonMarshal(&payload{
			Content: replace(bob.String(), repl),
		}).Post(notifier.DiscordWebhookURL)
		return err
	}
}
