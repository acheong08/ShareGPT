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

Response: A float64 with the amount of credit remaining

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
