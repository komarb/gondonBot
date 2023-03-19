## GOndon Discord Bot

A simple bot written in Go for Discord enjoyment purposes

## Features

- Music player which is controlled by buttons and not text commands (very cool)
- Audio modes: `normal`, `sped-up`, `slowed+reverb`, `bassboost`
- Meme (demotivator) generation based on messages and attachments (imgs, gifs) from your Discord server

## Usage

- `/help` - Displays list of all available commands
- `/play [youtube URL|youtube search query]` - Searches YouTube video/playlist by URL/search query and adds it to queue
- `/player` - Makes bot send message with embedded music player
- `/pop` - Remove last song from a queue
- `/mem` - Generate meme (demotivator) with random messages and attachment from server. Reply to a message with attachment with this command to make meme with specified attachment.
- `/popusk` - Finds out what user is a loser today

## Requirements

- Docker
- MondoDB

## Installation

Install Docker and in project root directory run this command:
```
docker run -d --env-file <path/to/config/env/file> --name <your_name> <your_name>
```

## Configuration

Config .env file must contain next content:
```
GDN_PREFIX=*                                       # prefix symbol for a command message
GDN_SERVICE_URL=https://www.googleapis.com         # url for youtube/tenor api requests
GDN_BOT_TOKEN=123YOURTOKEN123                      # your discord bot token
GDN_GAME_STATUS="gondon bot best bot"              # server bot game status
GDN_GOOGLE_API_KEY=123YOURAPIKEY123                # google api key
GDN_DB_URI=mongodb://localhost:27017/Gondon        # mongodb uri
```