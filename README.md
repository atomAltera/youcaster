# Youcaster

Youcaster is a simple tool to create a podcast episodes from a YouTube videos.

## Usage

First register telegram bot and get the token.
Also, you need to generate a Google API key with access to YouTube Data API (used to fetch video details like description, duration etc.)

You can restrict access to bot by setting `TELEGRAM_CHATS` env var. To get your chat ID, use [userinfobot](https://t.me/userinfobot)

Deploy Youcaster on your server. Here is an example docker-compose.yml file:

```yml
version: '2'

services:
  mongo:
    image: mongo:6
    restart: unless-stopped
    volumes:
      - ./mongo/db:/data/db
      - ./mongo/configdb:/data/configdb

  youcaster:
    image: atomaltera/youcaster:latest
    restart: unless-stopped
    environment:
      PUBLIC_BASE_URL: 'https://youcaster.example.com' # URL of your server
      DOWNLOAD_PATH: "/files" 
      MONGO_URI: 'mongodb://mongo/youcaster'
      WEB_ADDR: '0.0.0.0:3000'
      GOOGLE_API_KEY: '<google API key with access to YouTube Data API'
      TELEGRAM_CHATS: '<your tg chat id>'
      TELEGRAM_TOKEN: '<yout tg bot token>'
    ports:
      - '80:3000'
    depends_on:
      - mongo
    volumes:
      - './files:/files' # downloaded episodes
```


If your domain name is youcaster.example.com, feed URL will be http://youcaster.example.com/feed

Add it to your podcast player, send link to YouTube video to your telegram bot and enjoy!