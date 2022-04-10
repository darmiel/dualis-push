package notifier

import (
	"errors"
	"github.com/darmiel/dualis-push/dualis"
	"github.com/gregdel/pushover"
)

var (
	ErrPushoverTokenMissing     = errors.New("pushover token missing")
	ErrPushoverRecipientMissing = errors.New("pushover recipient missing")
)

func init() {
	notifiers["pushover"] = func(notifier *Notifier, grade *dualis.Grade, repl Replacers) error {
		if notifier.PushoverToken == "" {
			return ErrPushoverTokenMissing
		}
		if notifier.PushoverRecipient == "" {
			return ErrPushoverRecipientMissing
		}

		c := pushover.New(notifier.PushoverToken)
		r := pushover.NewRecipient(notifier.PushoverRecipient)
		f := notifier.Formatting()

		_, err := c.SendMessage(
			pushover.NewMessageWithTitle(
				replace(f.NewGradeMessageBody, repl),
				replace(f.NewGradeMessageTitle, repl),
			),
			r,
		)
		return err
	}
}
