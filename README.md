# twitchforwarder

Forwards Twitch chat to a Discord channel. 

Note that this uses Discord's webhooks, so this might break if the chat is particularly heavy. Caveat emptor.

[Downloads here.](https://github.com/strangebroadcasts/twitchforwarder/releases)

## Setup
1. Create a Discord webhook by opening the channel's properties, navigating to *Webhooks* and clicking *Create Webhook*. Make note of the webhook URL.
2. Optional: create your own OAuth token through [dev.twitch.tv](https://dev.twitch.tv/dashboard/apps/create). If you skip this step, twitchforwarder will automatically open up [a utility](https://twitchapps.com/tmi/) to generate this token.
3. Start twitchforwarder (through the command line for now:)

```
$ twitchforwarder -hook https://discordapp.com/api/webhooks/... -nick TwitchUser -channel VideoGameStream23 
```

where *nick* is your Twitch username, *hook* is the URL for the Discord webhook, and *channel* is the Twitch chat to forward messages from. If you created your own OAuth token, pass it with *-oauth*.

## Credits
This utility makes use of [gempir's](https://github.com/gempir) [go-twitch-irc](https://github.com/gempir/go-twitch-irc) and [pkg's](https://github.com/pkg) [browser](https://github.com/pkg/browser).
