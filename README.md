# `flaggio-cleaner-bot`

<p align="center">
  <img src="https://images.unsplash.com/photo-1563207153-f403bf289096?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&h=500&q=80">
</p>


Checks old Flaggio items in order to:
- Open an issue on Linear
- Notify to Slack

* You can configure it to do both or only one of the above items

## Install

### MacOS

Use `brew` to install it

```
brew tap brunoluiz/tap
brew install flaggio-cleaner-bot
```

### Linux and Windows

[Check the releases section](https://github.com/brunoluiz/flaggio-cleaner-bot/releases) for more information details 

### go get

Install using `GO111MODULES=off go get github.com/brunoluiz/flaggio-cleaner-bot/cmd/flaggio-cleaner-bot` to get the latest version. This will place it in your `$GOPATH`, enabling it to be used anywhere in the system.

**⚠️ Reminder**: the command above download the contents of master, which might not be the stable version. [Check the releases](https://github.com/brunoluiz/flaggio-cleaner-bot/releases) and get a specific tag for stable versions.

### Docker

The tool is available as a Docker image as well. Please refer to [Docker Hub page](https://hub.docker.com/r/brunoluiz/flaggio-cleaner-bot/tags) to pick a release

## Usage

The methods below can be used together

### Slack

1. Install in the Slack workspace and on the channel
1. `/invite @flaggio-cleaner-bot` or whatever name you had for the bot
1. Run the cli with the `--slack-token` and `--slack-channel` -- you should get a token once the application is added in the workspace

### Linear

1. Get a development token on Linear and run the cli using `--linear-token`, `--linear-project` and `--linear-team`

## Available params

```
NAME:
   flaggio-cleaner-bot - Process and notify old Flaggio feature flags 🤖

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --max-age value          Max flag age before processing is triggered (default: 120h0m0s) [$MAX_AGE]
   --storage-path value     Used to store already processed flags [$STORAGE_PATH]
   --flaggio-url value       [$FLAGGIO_URL]
   --flag-prefix value      Searches for a certain flag prefix -- useful for shared flaggio instances [$FLAG_PREFIX]
   --linear-token value      [$LINEAR_TOKEN]
   --linear-team value       [$LINEAR_TEAM]
   --linear-project value    [$LINEAR_PROJECT]
   --linear-template value  (default: "Clean-up code for feature flag '{{ .Key }}'") [$LINEAR_TEMPLATE]
   --slack-token value       [$SLACK_TOKEN]
   --slack-channel value     [$SLACK_CHANNEL]
   --slack-template value   (default: "Time to clean-up code for feature flag `{{ .Key }}` 🔥") [$SLACK_TEMPLATE]
   --help, -h               show help (default: false)
```
