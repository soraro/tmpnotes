# TMPNOTES
[![.github/workflows/ci.yml](https://github.com/soraro/tmpnotes/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/soraro/tmpnotes/actions/workflows/ci.yml)
[https://tmpnotes.com](https://tmpnotes.com)


Temporary notes that disappear after reading

# How it works
TMPNOTES is a system to share secrets in a reasonably secure, temporary way. TMPNOTES uses a [redis](https://redis.io/) cache to store notes. Notes are immediately purged once a note has been read, or the user-defined TTL (time to live) has expired. Notes are only allowed a maximum TTL of 24 hours. All notes are encrypted [server side](docs/server-side-encryption.md), and only the user who generated the note has the full key.

# Privacy
All notes are encrypted using our [server side encryption feature](docs/server-side-encryption.md) so we can not, and will not ever attempt to read notes stored in the redis database. Since these notes are meant to be temporary, no database backups are ever taken for [tmpnotes.com](https://tmpnotes.com). In addition, the website has a feature that can encrypt your note in the browser *before* it is sent to the server. These measures ensure that even if someone did view the database, the data would not be readable. The codebase is also small, and relatively easy to audit.

If you don't want to send your information to [tmpnotes.com](https://tmpnotes.com), we don't blame you! We want to make this project easy to run yourself if that is desired. We will host [tmpnotes.com](https://tmpnotes.com) as long as it does not become frequently abused or prohibitively expensive for us to do so.

# Run it yourself
We understand that secrets are sensitive and people may not want to use a publicly hosted instance ([tmpnotes.com](https://tmpnotes.com)). The following is information you can use to run a TMPNOTES system yourself! We strongly suggest that any implementation you run yourself is fronted by a reverse proxy (such as NGINX) with TLS.

## Environment variables
The following table shows environment variables that can be used to configure your TMPNOTES installation:
| Env var | Type | Description | Default |
|---------|------|-------------|---------|
| `TMPNOTES_ENABLE_HSTS` | bool | Return `"Strict-Transport-Security", "max-age=15552000"` header to enforce TLS for web browser clients. Only use this if you are sure your instance is running behind a reverse proxy with TLS. | `false` |
| `TMPNOTES_MAX_EXPIRE` | int | The maximum number of hours allowed before a note expires. | `24` |
| `TMPNOTES_MAX_LENGTH` | int | The maximum length (in characters) that a note is allowed to be. This should always be larger than `TMPNOTES_UI_MAX_LENGTH` to give room for the optional encryption padding in the UI. | `1000` |
| `TMPNOTES_UI_MAX_LENGTH` | int | The maximum length (in characters) that a note is allowed to be in the UI. This value should always be less than `TMPNOTES_MAX_LENGTH` to give room for the optional encryption padding in the UI. | `512` |
| `TMPNOTES_PORT` | int | Port number for the application to use. The env var `PORT` can also be used. | `5000` |
| `TMPNOTES_REDIS_URL` | string | Redis URI / connection string. `REDIS_URL` can also be used. | `redis://localhost:6379` |

## docker-compose
We have provided a `docker-compose` file to easily build and host a functional TMPNOTES system.
```
docker-compose up
```
Navigate to [localhost:5000](http://localhost:5000) and you will see your own tmpnotes instance running.

## Docker image
```
docker pull ghcr.io/soraro/tmpnotes:latest
```
View all versions on the [ghcr package page](https://github.com/soraro/tmpnotes/pkgs/container/tmpnotes)

## Helm Chart
*Coming soon!*

# Text clients
We made TMPNOTES so it can easily be used for text clients such as `curl` / `wget`. The following example shows how you can use curl to send and receive a note on [tmpnotes.com](https://tmpnotes.com)
```
MESSAGE="test message!"
EXPIRE=1

# Send a note:
ID=$(curl -X POST -d "{\"message\":\"$MESSAGE\",\"expire\":$EXPIRE}" https://tmpnotes.com/new)

# Receive the note:
curl https://tmpnotes.com/id/$ID
```

# Build / Test
You will need golang version 1.17. In addition, docker is recommended to run a local redis instance.

```
docker run -d --rm -p 6379:6379 redis
go build .
./tmpnotes

# navigate to http://localhost:5000 to see your local tmpnotes instance
```

Test:
```
go test -v ./...
```

# Questions / Comments / Suggestions
Please feel free to [open an issue](https://github.com/soraro/tmpnotes/issues/new)!


