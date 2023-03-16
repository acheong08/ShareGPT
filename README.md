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
  "object": "credit_summary",
  "total_granted": 18,
  "total_used": 1.850212,
  "total_available": 16.149788,
  "grants": {
    "object": "list",
    "data": [
      {
        "object": "credit_grant",
        "id": "f7855bd6-d87d-4ab6-8a6b-1c1b5b5b5b5b",
        "grant_amount": 18,
        "used_amount": 1.850212,
        "effective_at": 1670976000,
        "expires_at": 1680307200
      }
    ]
  },
  "error": {
    "message": "",
    "type": "",
    "param": "",
    "code": ""
  }
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
