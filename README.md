# ShareGPT

To allow the sharing of API keys to create a free OpenAI API.

## Requirements

- [Go](https://golang.org/)
- [Redis](https://redis.com/)

## Configuration

```bash
export REDIS_ADDRESS="HOST:PORT"
export REDIS_PASSWORD="..."
```

## Installation

`go install github.com/acheong08/ShareGPT@latest`

`export PATH=$PATH:$(go env GOPATH)/bin`

## Running

`ShareGPT`

## API

### GET /ping

Response: `{"message": "pong"}`

### POST /api_key/submit

Request:

```json
{ "api_key": "..." }
```

Response:

```json
{
  "object": "billing_subscription",
  "has_payment_method": true,
  "canceled": false,
  "canceled_at": null,
  "delinquent": null,
  "access_until": 1630454400,
  "soft_limit": 1600000,
  "hard_limit": 166666666,
  "system_hard_limit": 166666666,
  "soft_limit_usd": 96,
  "hard_limit_usd": 9999.99996,
  "system_hard_limit_usd": 9999.99996,
  "plan": {
    "title": "Pay-as-you-go",
    "id": "payg"
  },
  "account_name": "",
  "po_number": null,
  "billing_email": null,
  "tax_ids": null,
  "billing_address": {
    "city": "Penang",
    "line1": "<REDACTED>",
    "line2": null,
    "country": "MY",
    "postal_code": "<REDACTED>"
  },
  "business_address": null
}

```

### POST /api_key/delete

Request:

```json
{ "api_key": "..." }
```

Response:

```json
{
  "message": "API key deleted"
}
```

### POST /v1/chat

**This is the same as OpenAI's API**

```bash
curl http://HOST:PORT/v1/chat \
 -H 'Content-Type: application/json' \
 -d '{
  "model": "gpt-3.5-turbo",
  "messages": [
    {
      "role": "user",
      "content": "Say this is a test"
    }
  ]
}'
```
