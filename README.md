# dualis-push

(contains a [Dualis](https://dualis.dhbw.de/) "API" for Go)

---

Small service which can send a push message should a new grade be entered using [Pushover](http://pushover.net).
Discord support coming soon? 👀

## Example Runner

Multiple runners can be entered in `runners.toml` using the following format:

```toml
[[runners]]
user = "daniel@dhbw.de"
password = "kBsSduxLminWUhaCjjC7Rwnaifz8jncL"
cron = "0 */5 * * * *" # every 5m
PushoverToken = "cu3sm42ms449b9etcgwtiv3h4gfu6sar"
PushoverRecipient = "grmwfqr5cb7bkjdas7vpfq83bk9p4vz5"

[[runners]]
user = "brown@dhbw.de"
password = "msxMeMzTxf6ELu5r454w38X2khCCZLvM"
cron = "0 */30 * * * *" # every 30m
PushoverToken = "imuztspjkqyjzs4u3j34hn9nhdu5qm4r"
PushoverRecipient = "nd73gwqcp64fjqfmfg5xmfd5ygbm3tac"
```
