# The complete setup process

## Redis

This repository uses Redis to store the API key. The most popular provider of a Redis server is [redis.com](https://redis.com/). You can sign up for an account there for free.

- Take the public endpoint: e.g. `redis-11143.c8.us-west-1-2.ec1.cloud.redislabs.com:11143`
- Take the default user password

Save those to environment variables in a `.env` file

```bash
export REDIS_ADDRESS="HOST:PORT"
export REDIS_PASSWORD="..."
```

## Hosting

You'll probably want a VPS for the best reliability and speed. I DigitalOcean, Microsoft Azure, and Oracle Cloud are all great options.

SSH into your machine and install `git` and `golang` via `apt` from the root user.

```bash
$ sudo apt install git golang -y
```

Then to get this repository:
```bash
$ git clone https://github.com/acheong08/ShareGPT

$ cd ShareGPT

$ go build
```

This will compile a binary with the filename `ShareGPT`.

## Running

You need to ensure the environment variables are loaded
```bash
$ source .env
```

and then to run
```bash
PORT=8080 ./ShareGPT
```
You can set PORT to anything

## Other configurations

### NGINX

With `nginx`, you need to add a few options to allow streaming to happen efficiently.

Here is my configuration:
```nginx
location /share/ {
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                client_max_body_size 0;
                proxy_request_buffering off;
                proxy_read_timeout 300;
                proxy_max_temp_file_size 0;
                proxy_buffering off;
                proxy_no_cache 1; proxy_cache_bypass 1;
                proxy_pass http://localhost:8082/;
        }
```
### Cloudflare

You can use Cloudflare to handle rate limiting etc. It should be intuitive so I won't post a full guide here.
