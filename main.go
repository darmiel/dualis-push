package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/darmiel/dualis-push/dualis"
	"github.com/gregdel/pushover"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
}

type Runner struct {
	User     string
	Password string
	Cron     string
	grades   dualis.Grades

	PushoverToken     string
	PushoverRecipient string
}

func main() {
	// read runners config
	var (
		data []byte
		err  error
	)
	data, err = os.ReadFile("runners.toml")
	if err != nil {
		log.WithError(err).Fatal("cannot read runners.toml")
		return
	}

	// parse runner config
	type config struct {
		API     string
		Runners []Runner
	}
	var cfg config
	if _, err = toml.Decode(string(data), &cfg); err != nil {
		log.WithError(err).Fatal("cannot decode runners")
		return
	}

	log.Infof("API: %s", cfg.API)

	c := cron.New(cron.WithSeconds())
	for _, r := range cfg.Runners {
		if _, err := c.AddFunc(r.Cron, r.run); err != nil {
			log.WithError(err).Fatalf("Cannot create cron func for user %s", r.User)
			return
		}
		log.Infof("Added runner %s [%s]", r.User, r.Cron)
	}

	log.Info("Started cron task. Press CTRL-C to abort.")

	go func() {
		c.Run()
	}()

	// wait for CTRL-C
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// cancel cron
	c.Stop()
}

func (r *Runner) SendGradeUpdate(g *dualis.Grade) {
	log.Debugf("[%s] Got Updated Grade: %+v", g)

	if r.PushoverToken != "" && r.PushoverRecipient != "" {
		c := pushover.New(r.PushoverToken)
		res, err := c.SendMessage(
			pushover.NewMessageWithTitle(
				fmt.Sprintf("🎉 In [%s] hast du %s Pommes erhalten.", g.CourseName, g.Grade),
				"🫣 New Grade arrived, fat fuck!",
			),
			pushover.NewRecipient(r.PushoverRecipient),
		)
		if err != nil {
			log.WithError(err).Warn("Cannot send pushover message")
			return
		}
		log.Infof("Sent Pushover message: %+v", res)
	}
}

func (r *Runner) run() {
	log.Debugf("[%s] Checking user ...", r.User)

	client, err := dualis.Login(r.User, r.Password)
	if err != nil {
		log.WithError(err).Warn("cannot login to user (bad username/password?)")
		return
	}

	grades, err := client.CourseResults()
	if err != nil {
		log.WithError(err).Warn("cannot request grades")
		return
	}

	log.Debugf("[%s] Comparing old with new grades and looking for updates ...", r.User)
	r.grades.CheckForNewIn(grades, r.SendGradeUpdate)
	r.grades = grades // update new grades
}
