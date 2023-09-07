# Inspiration Bot

Uses the [InspiroBot](http://inspirobot.me/) API to generate random inspirational quotes and images.

Note: This bot is not affiliated with the creators of InspiroBot.

## Features

This bot adds 2 commands to your server:

- `/inspire` - sends a machine-generated inspirational quote and image
- `/source` - links to the source code of this bot

# Usage

1. Invite the bot to your server <https://discord.com/api/oauth2/authorize?client_id=1033443886472904764&permissions=0&scope=bot%20applications.commands>

2. Call `/inspire` !

I do keep some metrics of the bot usage, feel free to self-host it instead!

## Development

`DISCORD_TOKEN` is required as environment variables. These are set in a `.env` file like so:

```text
DISCORD_TOKEN=<token>
```

To learn more about discord bot development, visit [discord developers docs](https://discord.com/developers/docs/intro).