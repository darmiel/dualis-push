package discord

import (
	"bytes"
	"fmt"
	"github.com/apex/log"
	"github.com/darmiel/dualis-push/dualis"
	"github.com/imroc/req/v3"
	"strings"
	"text/template"
)

func Send(webhookURL, userID, userAvatarHash string, grades dualis.Grades) error {
	tpl := template.New("discord")
	tpl.Funcs(template.FuncMap{
		"Unit": func(g *dualis.Grade) string {
			return g.Unit()
		},
		"Replace": func(in string, repl ...string) string {
			res := in
			for i := 0; i < len(repl); i += 2 {
				res = strings.Replace(res, repl[i], repl[i+1], -1)
			}
			return res
		},
	})

	type data struct {
		Mention        string
		UserAvatarURL  string
		UserAvatarHash string
		Semester       string
		UserID         string
		WebhookURL     string
		Grade          *dualis.Grade
		Grades         dualis.Grades
	}

	d := &data{
		UserAvatarHash: userAvatarHash,
		UserID:         userID,
		WebhookURL:     webhookURL,

		Mention:       fmt.Sprintf("<@%s>", userID),
		UserAvatarURL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=512", userID, userAvatarHash),

		Grades: grades,
	}

	exec := func(input string) string {
		t, err := tpl.Parse(input)
		if err != nil {
			log.WithError(err).Warn("cannot parse template")
			return ""
		}
		var buf bytes.Buffer
		if err = t.Execute(&buf, d); err != nil {
			log.WithError(err).Warn("cannot execute template")
			return ""
		}
		return strings.TrimSpace(strings.Trim(buf.String(), "\n"))
	}

	p := DefaultPayload
	p.Content = exec(p.Content)
	p.Username = exec(p.Username)
	p.AvatarURL = exec(p.AvatarURL)

	// order grades by semester
	semesters := make(map[string][]*dualis.Grade)
	for _, g := range grades {
		semesters[g.Semester] = append(semesters[g.Semester], g)
	}

	for s, gr := range semesters {
		d.Semester = s // update template data

		e := DefaultEmbed // create new default embed by discord config
		e.Title = exec(e.Title)
		e.Description = exec(e.Description)
		e.Thumbnail.URL = exec(e.Thumbnail.URL)
		e.Author.IconUrl = exec(e.Author.IconUrl)
		e.Author.Name = exec(e.Author.Name)
		e.Author.Url = exec(e.Author.Url)
		e.Footer.Text = exec(e.Footer.Text)

		for _, g := range gr {
			d.Grade = g // update template data

			f := DefaultField
			f.Name = exec(f.Name)
			f.Value = exec(f.Value)

			e.Fields = append(e.Fields, f)
		}

		p.Embeds = append(p.Embeds, e)
	}

	r, err := req.R().SetBodyJsonMarshal(p).Post(webhookURL)
	log.Debugf("Webhook Status: %d", r.StatusCode)
	log.Debugf("Webhook Response: %s", r.String())

	return err
}
