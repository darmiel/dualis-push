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
	config    *config
}

// parse runner config
type config struct {
	NotifyFirst bool
	Runners     []Runner
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

	var cfg config
	if _, err = toml.Decode(string(data), &cfg); err != nil {
		log.WithError(err).Fatal("cannot decode runners")
		return
	}

	c := cron.New(cron.WithSeconds())
	for _, r := range cfg.Runners {
		co := r

		co.grades = make(dualis.Grades)
		co.config = &cfg
		if _, err := c.AddFunc(co.Cron, func() {
			co.run()
		}); err != nil {
			log.WithError(err).Fatalf("Cannot create cron func for user %s", co.User)
			return
		}
		log.Infof("Added runner %s [%s]", co.User, co.Cron)
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

	// collect new grades
	updates := make(dualis.Grades)
	for m, g := range grades {
		if _, ok := r.grades[m]; !ok {
			r.grades[m] = g
			updates[m] = g
		}
	}

	if len(updates) > 0 && (r.config.NotifyFirst || r.check) { // ignore first fetch
		log.Debugf("[%s] Sending %d grade updates...", r.User, len(updates))
		r.sendGradeUpdate(updates)
	}

	r.check = true // mark first fetch
}

func (r *Runner) sendGradeUpdate(g dualis.Grades) {
	log.Debugf("[%s] Got Updated Grade: %+v", r.User, g)

	for _, n := range r.Notifiers {
		if n.Disabled {
			continue
		}
		log.Debugf("[notify@%s] sending grade", n.Type)
		if err := n.DoGradesArrived(g); err != nil {
			log.WithError(err).Warnf("[notify@%s]error sending notification", n.Type)
		}
	}
}
