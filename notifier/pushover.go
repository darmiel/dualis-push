package notifier

import (
	"errors"
	"github.com/darmiel/dualis-push/dualis"
	"github.com/darmiel/dualis-push/notifier/discord"
	"github.com/gregdel/pushover"
)

var (
	ErrPushoverTokenMissing     = errors.New("pushover token missing")
	ErrPushoverRecipientMissing = errors.New("pushover recipient missing")
)
var (
	ErrDiscordWebhookMissing        = errors.New("webhook url missing")
	ErrDiscordUserIDMissing         = errors.New("user id missing")
	ErrDiscordUserAvatarHashMissing = errors.New("user avatar hash missing")
)

func init() {
	notifiers["pushover"] = func(notifier *Notifier, grades dualis.Grades) (err error) {
		if notifier.PushoverToken == "" {
			return ErrPushoverTokenMissing
		}
		if notifier.PushoverRecipient == "" {
			return ErrPushoverRecipientMissing
		}

		c := pushover.New(notifier.PushoverToken)
		r := pushover.NewRecipient(notifier.PushoverRecipient)
		f := notifier.Formatting()

		for _, g := range grades {
			if _, err = c.SendMessage(
				pushover.NewMessageWithTitle(
					replaceGrade(f.NewGradeMessageBody, g),
					replaceGrade(f.NewGradeMessageTitle, g),
				),
				r,
			); err != nil {
				return
			}
		}
		return
	}
	notifiers["discord"] = func(notifier *Notifier, grades dualis.Grades) error {
		if notifier.DiscordWebhookURL == "" {
			return ErrDiscordWebhookMissing
		}
		if notifier.DiscordUserID == "" {
			return ErrDiscordUserIDMissing
		}
		if notifier.DiscordAvatarHash == "" {
			return ErrDiscordUserAvatarHashMissing
		}
		return discord.Send(notifier.DiscordWebhookURL, notifier.DiscordUserID, notifier.DiscordAvatarHash, grades)
	}
}
