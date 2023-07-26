[![Carousell Gobot Logo](carousell-gobot.png)](https://gitlab.com/jarylc/carousell-gobot)

# Carousell GoBot
A Carousell.sg automation bot written in Golang
> WARNING: NIGHTLY RELEASES ONLY. Expect bugs and possibly breaking changes every update.

[**Docker Hub Image »**](https://hub.docker.com/r/jarylc/carousell-gobot)

[**Explore the docs »**](https://gitlab.com/jarylc/carousell-gobot)

[Report Bugs](https://gitlab.com/jarylc/carousell-gobot/-/issues/new?issuable_template=Bug)
· [Request Features](https://gitlab.com/jarylc/carousell-gobot/-/issues/new?issuable_template=Feature%20Request)


## About The Project
### Demo
#### Initial auto-reply
![Initial auto-reply](readme/demo-initial.gif)

#### Low-balling detection
![Low-balling detection](readme/demo-lowball.gif)

#### Lower official offer detection
![Lower official offer detection](readme/demo-loweroffer.gif)

#### Forwarding to other chat services
![Forwarding to other chat services](readme/demo-forwarding.jpg)

#### Reminders & Commands
![Reminders & Commands](readme/demo-reminder.gif) ![Reminders](readme/demo-reminder.jpg)


### Built With
* [golang](https://golang.org/)
* [gorrilla/websocket](https://github.com/gorilla/websocket)
* [dlclark/regexp2](https://github.com/dlclark/regexp2)
* [dop251/goja](https://github.com/dop251/goja)
  * [jarylc/go-chrono](https://github.com/jarylc/go-chrono)
      * [wanasit/chrono](https://github.com/wanasit/chrono)
* [jarylc/go-chromedpproxy](https://github.com/jarylc/go-chromedpproxy)


## Getting Started
To get a local copy up and running follow these simple steps.

### Docker (recommended)
#### 1. Prepare configuration file
1. [Click here to download the latest sample config](https://gitlab.com/jarylc/carousell-gobot/-/raw/master/config.sample.yaml?inline=false)
2. Rename `config.sample.yaml` to `config.yaml`
3. Edit `config.yaml` with your configurations (breakdown below)

#### 2a. Docker Run
```shell
docker run -it -d \
  --name carousell-gobot \
  -e TZ=Asia/Singapore \
  -e INSTALL_CHROME=1 `# OPTIONAL: install Chromium before start for username & password login support` \
  -p 9221:9221 `# OPTIONAL: configured forward port for front-end portal (for 2FA & Captcha entries)` \
  -v /path/to/config/folder:/data \
  jarylc/carousell-gobot
```

#### 2b. Docker-compose
```docker-compose
carousell-gobot:
    image: jarylc/carousell-gobot
    environment:
      - TZ=Asia/Singapore
      - INSTALL_CHROME=1 # OPTIONAL: install Chromium before start for username & password login support
    ports:
      - "9221:9221" # OPTIONAL: configured forward port for front-end portal (for 2FA & Captcha entries)
    volumes:
      - /path/to/config/folder:/data
```

### Binaries
1. [Visit here to view and download the latest build artifacts](https://gitlab.com/jarylc/carousell-gobot/-/jobs/artifacts/master/browse?job=build) / Alternatively, build it yourself with instructions below
2. [Click here to download the latest sample config](https://gitlab.com/jarylc/carousell-gobot/-/raw/master/config.sample.yaml?inline=false)
3. Rename `config.sample.yaml` to `config.yaml`
4. Edit `config.yaml` with your configurations (breakdown below)
5. Run the binary
> Note: you can also specify path to config using `-c /path/to/config.yaml` and path to state using `-s /path/to/state.json`


## Development
### Building native binary
```shell
$ cd /path/to/project/folder
$ go build -ldflags="-w -s"
```

### Building for 64-bit Windows
```shell
$ cd /path/to/project/folder
$ GOOS=windows GOARCH=amd64 go build -ldflags="-w -s"
```

### Docker build
```shell
$ cd /path/to/project/folder
$ docker build .
```

## Usage
### If username and password not specified, get Carousell.sg Cookie and Token
#### Firefox
1. Go to https://www.carousell.sg/ds/api/1.0/chat/token/?_path=/1.0/chat/token/
2. Press `CTRL+SHIFT+E` / `⌘+⌥+E`
3. Refresh the page
4. ![Firefox Final Step](readme/cookie-firefox.png)
5. Copy the value of `token` for token as well

#### Chrome
1. Go to https://www.carousell.sg/ds/api/1.0/chat/token/?_path=/1.0/chat/token/
2. Press `CTRL+SHIFT+C` / `⌘+⌥+C`
3. Click on `Network` tab
4. Refresh the page
5. ![Chrome Final Step](readme/cookie-chrome.png)
6. Copy the value of `token` for token as well

### Commands (prefixed with command_prefix configuration)
Currently, only your account can send these commands:
- `sched`/`schedule`/`remind`/`reminder`/`deal` - schedule a deal and set reminders based on reminders configuration
  - arguments: date/day and time (optional)
  - if no arguments, it will read the most recent reply and response
  - uses natural date processing. ELI5: something like a simple AI to detect dates in sentences.
- `cancel`/`del`/`delete` - cancel deal & reminders
- `faq` - resend the FAQ
- `contact` - send contact details
- `stop` - immediately stop the bot

### Configuration
#### Sample
https://gitlab.com/jarylc/carousell-gobot/-/blob/master/config.sample.yaml

#### Breakdown
- `application` - application related configurations
  - `base_url` - base URL of front-end portal (to be sent to forwarders on 2FA or Captcha requests)
  - `portal_listener` - listen address for front-end
  - `chrome_listener` - listen address for Chrome remote debugger
- `carousell` - Carousell related configurations
  - `cookie` - either this and `token` or `username` and `password` or both are required, entire value inside Cookie header from above
  - `token` - either this and `cookie` or `username` and `password` or both are required, retrieved from https://www.carousell.sg/ds/api/1.0/chat/token/?_path=/1.0/chat/token/ after logging in
  - `username` - either this or `cookie` and `token` or both are required, specify username for automatic session renewal (if Docker, make sure to run with env `INSTALL_CHROME=1`)
  - `password` - either this or `cookie` and `token` or both are required, specify password for automatic session renewal (if Docker, make sure to run with env `INSTALL_CHROME=1`)
  - `ping_interval` - interval to ping Carousell to check connectivity
  - `low_ball` - percentage of price to be considered low-ball in decimal notation
- `message_templates` - message templates
  - `faq` - frequently answered questions (optional, leave empty if not required)
  - `initial` - new chat initial greeting (placeholders available)
    - `{{NAME}}` - name of sender
    - `{{ITEM}}` - name of item
  - `not_available` - chat message on a listing that is sold, reserved or deleted (placeholders available)
    - `{{REASON}}` - reason why listing is not available anymore
  - `offered` - official offer using Carousell's offer button (placeholders available)
    - `{{OFFER}}` - offer value
  - `possible_offer` - offer value detected from previously sent message (placeholders available)
    - `{{OFFER}}` - offer value
  - `lower_offer` - official offer value was decreased from possible offer, buyer being cheeky?
  - `low_balled` - low-balled (placeholders available)
    - `{{PERCENT}}` - low-ball percent configured in configuration in percent notation
  - `contact` - contact details reply for contact command
  - `reminder` - reminder (placeholders available)
    - `{{HOURS}}' - hours remaining till scheduled time
- `reminders` - list of hours to set reminders for on every deal (optional, leave empty if none)
- `command_prefix` - prefix to use for the app to detect as a command
- `state_prune` - number of days to keep items in state after last activity
- `forwarders` - list of forwarders (optional, `[]` if none)
  - for `telegram` (guide: https://core.telegram.org/bots):
    - ```yaml
      - type: telegram
        token: '1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZabcde1234' # bot token
        chat_id: 87654321 # chat id
        message_templates:
        standard: |- # standard message
          *Carousell Chat/Offer*
          [{{NAME}} - {{ITEM}}](https://www.carousell.sg/inbox/{{ID}})
          ${{OFFER}}
          {{FLAGS}}
        reminder: |- # reminder message
          *Carousell Deal Reminder*
          [{{ITEM}}](https://www.carousell.sg/inbox/{{ID}})
          Deal ${{OFFER}} in {{HOURS}} hour(s)!
      ```
  - for `discord` (webhook only, guide: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks):
    - ```yaml
      - type: discord
        webhook_url: 'https://discord.com/api/webhooks/012345678910111213/lrennsXdocFAKEdOWAkpWEBHOOKMrMMcXa_FcmCURLUMNU_STsnbKVTRirmvLccmvkpA' # webhook url
        message_templates:
        standard: |- # standard message
          **Carousell Chat/Offer**
          [{{NAME}} - {{ITEM}}](https://www.carousell.sg/inbox/{{ID}})
          ${{OFFER}}
          {{FLAGS}}
        reminder: |- # reminder message
          **Carousell Reminder**
          [{{ITEM}}](https://www.carousell.sg/inbox/{{ID}})
          Deal ${{OFFER}} in {{HOURS}} hour(s)!
      ```
  - for `slack` (webhook only, guide: https://api.slack.com/messaging/webhooks):
    - ```yaml
      - type: slack
        webhook_url: 'https://hooks.slack.com/services/T1234GN5678/B1FAKEZNRCM/fPoWEBHOOK2Yyn3oURLOjW4J' # webhook url
        message_templates:
        standard: |- # standard message
          *Carousell Chat/Offer*
          <https://www.carousell.sg/inbox/{{ID}}|{{NAME}} - {{ITEM}}>
          ${{OFFER}}
          {{FLAGS}}
        reminder: |- # reminder message
          *Carousell Reminder*
          <https://www.carousell.sg/inbox/{{ID}}|{{ITEM}}>
          Deal ${{OFFER}} in {{HOURS}} hour(s)!
      ```
      

## Roadmap
See the [open issues](https://gitlab.com/jarylc/carousell-gobot/-/issues) for a list of proposed features (and known
issues).


## Contributing
Feel free to fork the repository and submit pull requests.


## License
Distributed under the GNU GENERAL PUBLIC LICENSE V3. See `LICENSE` for more information.


## Contact
Jaryl Chng - git@jarylchng.com

https://jarylchng.com

Project Link: [https://gitlab.com/jarylc/carousell-gobot/](https://gitlab.com/jarylc/carousell-gobot/)
