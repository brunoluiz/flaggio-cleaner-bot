# `flaggio-cleaner-bot`

<p align="center">
  <img src="https://images.unsplash.com/photo-1563207153-f403bf289096?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&h=500&q=80">
</p>


Checks old Flaggio items in order to:
- Open an issue on Linear
- Notify to Slack

* You can configure it to do both or only one of the above items

## Installing

1. Install in the Slack workspace and on the channel
1. `/invite @flaggio-cleaner-bot` or whatever name you had for the bot

## Available params

```
   --max-age value          (default: 120h0m0s) [$MAX_AGE]
   --db-dsn value            [$DB_DSN]
   --flag-prefix value       [$FLAG_PREFIX]
   --linear-token value      [$LINEAR_TOKEN]
   --linear-team value       [$LINEAR_TEAM]
   --linear-project value    [$LINEAR_PROJECT]
   --linear-template value  (default: "Clean-up code for feature flag '{{ .Key }}'") [$LINEAR_TEMPLATE]
   --slack-token value       [$SLACK_TOKEN]
   --slack-channel value     [$SLACK_CHANNEL]
   --slack-template value   (default: "Time to clean-up code for feature flag `{{ .Key }}` 🔥") [$SLACK_TEMPLATE]
```

## Wishlist

- Replace Mongo access on Flaggio with a GraphQL access
- Evaluate another solution for dedupe notifications (using flaggio mongo db now) -- badger?
