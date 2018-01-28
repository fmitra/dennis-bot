# Dennis

A pet project to learn Go, Dennis is Telegram bot to manage expense tracking.

## Development

To get started you'll need a [Telegram auth token](https://core.telegram.org/bots/api#authorizing-your-bot) and [Ngrok](https://ngrok.com/download)

Telegram does not send any authentication headers in their requests, and instead recommends you instead use the token as the path of your webhook.

1. Set your Telegram token to your development environment

```
export TELEGRAM_AUTH_TOKEN=$YOUR_TOKEN
export WITAI_AUTH_TOKEN=$YOUR_TOKEN
```

2. Run the bot and ngrok

```
go build
./dennis
./ngrok http 8080
```

3. Send Telegram your webhook

```
curl --data "url=https://abcd.ngrok.io/<TELEGRAM_AUTH_TOKEN>" https://api.telegram.org/bot<TELEGRAM_AUTH_TOKEN>/setWebhook
```
