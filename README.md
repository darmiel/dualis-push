<img align="right" float="right" src="https://user-images.githubusercontent.com/71837281/137636053-f090fa9a-7a70-4028-a95b-89810263caac.png" height="256px" width="237px">

# dualis-push

Since my university doesn't send out notifications when new grades are entered, 
which is why you always have to manually check if a grade has been entered, 
I wrote this little service which fetches all grades every X minutes and 
then checks if new grades have been added.

If a new grade is found, the user is informed about the new grade via configurable notification channels.

## Notification Channels

### Discord Webhook

Requires `DiscordWebhookURL` in the notifier configuration

### Pushover

Requires `PushoverToken` and `PushoverRecipient` in the notifier configuration.

---

(contains a [Dualis](https://dualis.dhbw.de/) "API" for Go in the `dualis` package)

## Example Runner

Multiple runners can be entered in `runners.toml` using the following format:

```toml
[[runners]]
    User = "daniel@dh-karlsruhe.de"
    Password = "22z4o2ZoMkzqgf2m4cCmfrH5cye3XUML"
    Cron = "0 */10 * * * *" # check for new grades every 10 minutes

    [[runners.notifiers]]
        Type = "Pushover"
        Disabled = true # temporaily disable a notifier
        PushoverToken = "pgfx5j4w5dq9r2d95fr754xur24xniwd"
        PushoverRecipient = "c5aveuvgf47b3z2bptn3c33wunnmr8oe"

    [[runners.notifiers]]
        Type = "Discord"
        DiscordWebhookURL = "https://discord.com/api/webhooks/961920185726003913/s3inqyxkhm2jnsvf45o84mq-xt5i8v6ugkmzdd7hbc6unof65f_68z5seik9mcipuqx5"
        [runners.notifiers.format] # override formatting
            NewGradeMessageTitle = ""
            NewGradeMessageBody = """ \
            🚨 [@everyone] Neue Note wurde eingetragen\n \
            In **%course%** wurde vermutlich eine neue Note eingetragen!\n \
            👉 https://dualis.dhbw.de \
            """

[[runners]]
    User = "brown@dhbw.de"
    Password = "msxMeMzTxf6ELu5r454w38X2khCCZLvM"
    Cron = "0 */30 * * * *" # every 30 minutes
    [[runners.notifiers]]
        Type = "Pushover"
        PushoverToken = "imuztspjkqyjzs4u3j34hn9nhdu5qm4r"
        PushoverRecipient = "nd73gwqcp64fjqfmfg5xmfd5ygbm3tac"
```
