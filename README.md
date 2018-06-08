[![Build Status](https://travis-ci.org/fmitra/dennis-bot.svg?branch=master)](https://travis-ci.org/fmitra/dennis-bot)

# Dennis

A pet project to learn Go, Dennis is Telegram bot to manage expense tracking.

![dennis](https://www.francismitra.com/static/misc/dennis/convo.jpg)

## Overview

Dennis was written to track international expenses. He keeps a log of expenses in any
currency and returns the total (daily, weekly, monthly) in USD. At the moment he
supports the following commands:

* Track an expense

```
format: <integer_amount><currency_iso> for <description>

example: 200RUB for Lunch
```

* Get expense history

```
format: how much did I spend <time_period> (today, this week, this month)

example: How much did I spend today?
```

## Developer Dependencies

* [Ngrok](https://ngrok.com/downlaod)
* Postgres & Redis or [Docker](https://www.docker.com/)

## Getting Started

You will need API key's for the following services to get started.

* [Telegram Auth Token](https://core.telegram.org/bots/api#authorizing-your-bot)
* [Alphapoint API Key](https://www.alphapoint.com/api/index.html)
* [Wit.ai API Key](https://wit.ai)

#### 1. Set up development environment

The test suite will expect Postgres and Redis to be set up as well as a valid
configuraiton file. The `config.example.json` file is already prepared to use the
default settings in the sample `docker-compose.example.yml`.

```
make develop
docker-compose up -d
dep ensure -vendor-only -v
```

#### 2. Confirm tests are passing

```
go test ./...
```

#### 3. Run Ngrok

This step is not necessary if all you want to do is run the test suite.

```
./ngrok http 8080
```

#### 4. Set up your local `config.json`

* Postgres & Redis settings if you are not using the default test config
* Telegram API token to respond to messages
* Wit.ai auth token to parse user messages
* Alphapoint API key to convert currency
* Domain the bot will be receiving webhooks from. In development, this will be the Ngrok URL


#### 5. Run the bot

```
go build ./cmd/dennis-bot
./dennis-bot
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
