# Dennis

A pet project to learn Go, Dennis is Telegram bot to manage expense tracking.

## Development

To get started you'll need a [Telegram auth token](https://core.telegram.org/bots/api#authorizing-your-bot) and [Ngrok](https://ngrok.com/download)

Telegram does not send any authentication headers in their requests, and instead recommends you instead use the token as the path of your webhook.

### Config

Set up a `config.json` file using `config.example.json` as a template. The configuration file will require

* Postgres DB settings
* Telegram API token to respond to messages
* Wit.ai auth token to parse user messages
* Alphapoint API key to convert currency
* Domain the bot will be receiving webhooks from

### Docker set up

If Postgres and Redis are missing from the environment

```
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=dennis -e POSTGRES_USER=dennis -e POSTGRES_DB=dennis_test postgres
docker run -d --name redis -p 6379:6379 redis
```

Run Dennis (remove Postgres and Redis flags if they're running outside the container)

```
docker build -t fmitra/dennis .
docker run --rm -p 8080:8080 -v `pwd` fmitra/dennis --link=postgres:postgres --link=redis:redis
```

Start up Ngrok

```
./ngrok http 8080
```

### Non Docker setup

```
dep ensure -vendor-only -v
go build
./dennis
./ngrok http 8080
```
