[![Carousell Gobot Logo](carousell-gobot.png)](https://gitlab.com/jarylc/carousell-gobot)

# Carousell GoBot
A Carousell.sg automation bot written in Golang
> WARNING: IN HEAVY DEVELOPMENT AND ALPHA RELEASES ONLY. Expect bugs and possibly breaking changes every update.

[**Docker Hub Image »**](https://hub.docker.com/r/jarylc/carousell-gobot)

[**Explore the docs »**](https://gitlab.com/jarylc/carousell-gobot)

[Report Bugs](https://gitlab.com/jarylc/carousell-gobot/-/issues)
· [Request Features](https://gitlab.com/jarylc/carousell-gobot/-/issues)


## About The Project
### Demo
#### Initial auto-reply
![Initial auto-reply](readme/demo-initial.gif)
#### Low-balling detection
![Low-balling detection](readme/demo-lowball.gif)
#### Lower official offer detection
![Lower official offer detection](readme/demo-loweroffer.gif)
#### Reminders & Commands
![Reminders & Commands](readme/demo-reminder.jpg)
#### Forwarding to other chat services
![Forwarding to other chat services](readme/demo-forwarding.jpg)
### Built With
* [golang](https://golang.org/)
* [gorrilla/websocket](https://github.com/gorilla/websocket)
* [dlclark/regexp2](https://github.com/dlclark/regexp2)
* [dop251/goja](https://github.com/dop251/goja)
    * [wanasit/chrono](https://github.com/wanasit/chrono)


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
  -e TZ=UTC \
  -v /path/to/config/folder:/data \
  jarylc/carousell-gobot
```
#### 2b. Docker-compose
```docker-compose
carousell-gobot:
    image: jarylc/carousell-gobot
    environment:
      - TZ=UTC
    volumes:
      - /path/to/config/folder:/data
    restart: unless-stopped
```


### Binaries
1. [Visit here to view and download the latest build artifacts](https://gitlab.com/jarylc/carousell-gobot/-/jobs/artifacts/master/browse?job=build) / Alternatively, build it yourself with instructions below
2. [Click here to download the latest sample config](https://gitlab.com/jarylc/carousell-gobot/-/raw/master/config.sample.yaml?inline=false)
3. Rename `config.sample.json` to `config.yaml`
4. Edit `config.json` with your configurations (breakdown below)
5. Run the binary
> Note: you can also specify path to config using `-c /path/to/config.yaml` and path to state using `-s /path/to/state.json`


## Development
### Building native binary
```shell
$ cd /path/to/project/folder
$ go build -ldflags="-w -s"
```
### Building for amd64 Windows
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
### Get Carousell.sg Cookie
#### Firefox
1. Go to https://www.carousell.sg/
2. Right click anywhere
3. Press `CTRL+SHIFT+E`
4. Refresh the page
5. ![Firefox Final Step](readme/cookie-firefox.png)
### Configuration
#### Sample
https://gitlab.com/jarylc/carousell-gobot/-/blob/master/config.sample.yaml
#### Breakdown
- `carousell` - Carousell related configurations
  - `cookie` - entire value inside Cookie header from above
  - `ping_interval` - interval to ping Carousell to check connectivity
  - `low_ball` - percentage of price to be considered low-ball in decimal notation
- `message_templates` - message templates
  - `faq` - frequently answered questions (optional, leave empty if not required)
  - `initial` - new chat initial greeting (placeholders available)
    - `{{NAME}}` - name of sender
    - `{{ITEM}}` - name of item
  - `offered` - official offer using Carousell's offer button (placeholders available)
    - `{{OFFER}}` - offer value
  - `posssible_offer` - offer value detected from previously sent message (placeholders available)
    - `{{OFFER}}` - offer value
  - `lower_offer` - official offer value was decreased from possible offer, buyer being cheeky?
  - `low_balled` - low-balled (placeholders available)
    - `{{PERCENT}}` - low-ball percent configured in configuration in percent notation
  - `reminder` - reminder (placeholders available)
    - `{{HOURS}}' - hours remaining till scheduled time
- `reminders` - list of hours to set reminders for on every deal (optional, leave empty if none)
- `command_prefix` - prefix to use for the app to detect as a command
- `state_prune` - number of days to keep items in state after last activity
- `forwarders` - list of forwarders (optional, leave empty if none)
  - for `telegram`:
    - ```yaml
      - type: telegram
        token: 1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZabcde1234 # bot token
        chat_id: '87654321' # chat id
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