package discord

import (
	"github.com/BurntSushi/toml"
	"github.com/apex/log"
)

var (
	DefaultPayload discordPayload
	DefaultEmbed   embed
	DefaultField   field
)

func init() {
	type ft struct {
		Payload discordPayload
		Embed   embed
		Field   field
	}
	var f ft
	if _, err := toml.DecodeFile("discord.defaults.toml", &f); err != nil {
		log.WithError(err).Warn("cannot load discord defaults")
	} else {
		DefaultPayload = f.Payload
		DefaultEmbed = f.Embed
		DefaultField = f.Field
	}
}

type field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
type author struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"icon_url"`
}
type footer struct {
	Text string `json:"text"`
}

type thumbnail struct {
	URL string `json:"url"`
}

type embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Color       int       `json:"color"`
	Fields      []field   `json:"fields"`
	Thumbnail   thumbnail `json:"thumbnail"`
	Author      author    `json:"author"`
	Footer      footer    `json:"footer"`
}

type discordPayload struct {
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	Embeds    []embed `json:"embeds"`
}
