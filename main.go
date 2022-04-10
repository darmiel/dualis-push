package main

import (
	"github.com/BurntSushi/toml"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/darmiel/dualis-push/dualis"
	"github.com/darmiel/dualis-push/notifier"
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

	// used by checker
	grades dualis.Grades
	check  bool

	// notifiers
	Notifiers []*notifier.Notifier
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

	if r.check { // ignore first fetch
		log.Debugf("[%s] Comparing old with new grades and looking for updates ...", r.User)
		r.grades.CheckForNewIn(grades, r.sendGradeUpdate)
	}
	r.check = true
	r.grades = grades // update new grades
}

func (r *Runner) sendGradeUpdate(g *dualis.Grade) {
	log.Debugf("[%s] Got Updated Grade: %+v", r.User, g)

	for _, n := range r.Notifiers {
		if n.Disabled {
			continue
		}
		log.Debugf("[notify@%s] sending grade", n.Type)
		if err := n.DoGradeArrived(g); err != nil {
			log.WithError(err).Warnf("[notify@%s]error sending notification", n.Type)
		}
	}
}
