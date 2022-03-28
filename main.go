package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/gregdel/pushover"
	"github.com/imroc/req/v3"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
)

const (
	APIUrl = "http://localhost:5001"
)

func init() {
	log.SetHandler(cli.Default)
	log.SetLevel(log.DebugLevel)
}

type Runner struct {
	User     string
	Password string
	Cron     string
	grades   Grades

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
		Runners []Runner
	}
	var cfg config
	if _, err = toml.Decode(string(data), &cfg); err != nil {
		log.WithError(err).Fatal("cannot decode runners")
		return
	}

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

type Grade struct {
	Name               string `json:"name,omitempty"`
	Date               string `json:"date,omitempty"`
	Grade              string `json:"grade,omitempty"`
	ExternallyAccepted bool   `json:"externally accepted,omitempty"`
}

func (g *Grade) Marshal() (res string, err error) {
	var data []byte
	if data, err = json.Marshal(g); err != nil {
		return
	}
	res = string(data)
	return
}

type Grades []*Grade

func (g Grades) CheckForNewIn(other Grades, con func(g *Grade)) {
	for _, n := range other {
		// check if n in g
		nid, err := n.Marshal()
		if err != nil {
			log.Warn("cannot compare")
			return
		}

		found := false

	a:
		for _, o := range g {
			oid, err := o.Marshal()
			if err != nil {
				log.Warn("cannot compare")
				return
			}
			if oid == nid {
				found = true
				break a
			}
		}

		if !found {
			con(n) // send update
		}
	}
}

func (r *Runner) SendGradeUpdate(g *Grade) {
	log.Debugf("[%s] Got Updated Grade: %+v", g)

	if r.PushoverToken != "" && r.PushoverRecipient != "" {
		c := pushover.New(r.PushoverToken)
		res, err := c.SendMessage(
			pushover.NewMessageWithTitle(
				fmt.Sprintf("🎉 In [%s] hast du %s%% erhalten.", g.Name, g.Grade),
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

	var grades Grades
	res, err := req.R().
		SetResult(&grades).
		SetHeaders(map[string]string{
			"X-Auth-User": r.User,
			"X-Auth-Pass": r.Password,
		}).
		Get(APIUrl + "/dualis/api/v1.0/grades/")

	if err != nil {
		log.WithError(err).Warn("cannot request grades.")
		return
	}

	if res.StatusCode/100 != 2 {
		log.Warn("Status Code was not 2xx.")
		return
	}

	log.Debugf("[%s] Comparing old with new grades and looking for updates ...", r.User)
	r.grades.CheckForNewIn(grades, r.SendGradeUpdate)
	r.grades = grades // update new grades
}
