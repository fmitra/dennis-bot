# Dennis

A pet project to learn Go, Dennis is Telegram bot to manage expense tracking.

## Developer Dependencies

* [Ngrok](https://ngrok.com/downlaod)
* Postgres & Redis or [Docker](https://www.docker.com/)

## Getting Started

You will need API key's for the following services to get started.

* [Telegram Auth Token](https://core.telegram.org/bots/api#authorizing-your-bot)
* [Alphapoint API Key](https://www.alphapoint.com/api/index.html)
* [Wit.ai API Key](https://wit.ai)

#### Create your configuration files.

```
make develop
```

#### Edit `config.json` with the following settings

* Postgres DB settings
* Telegram API token to respond to messages
* Wit.ai auth token to parse user messages
* Alphapoint API key to convert currency
* Domain the bot will be receiving webhooks from. In development, this will be the Ngrok URL

#### Start Postgres and Redis

```
docker-compose up -d
```

#### Run Ngrok

```
./ngrok http 8080
```

#### Run the bot

```
dep ensure -vendor-only -v
go test ./...
go build
./dennis
```

## Developer Notes

#### Telegram Authentication

Telegram does not send any authentication headers in their requests, and instead recommends
you instead use the token as the path of your webhook.

#### Docker

There is a Dockerfile to build the bot for deploy. If you'd like to include it in with the
other dependencies in `docker-compose.yml` you can add the following to your file's services:

```
bot:
  build: .
  ports:
    - 8080:8080
  restart: unless-stopped
  depends_on:
    - postgres
    - redis
```
