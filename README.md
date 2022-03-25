# TMPNOTES
[![.github/workflows/ci.yml](https://github.com/soraro/tmpnotes/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/soraro/tmpnotes/actions/workflows/ci.yml)  
[https://tmpnotes.com](https://tmpnotes.com)


Temporary notes that disappear after reading

# How it works
TMPNOTES is a system to share secrets in a reasonably secure, temporary way. TMPNOTES uses a [redis](https://redis.io/) cache to store notes. Notes are immediately purged once a note has been read, or the user-defined TTL (time to live) has expired. Notes are only allowed a maximum TTL of 24 hours.

# Privacy
We will never read notes that are stored in the redis database, and there are no backups that are ever taken for [tmpnotes.com](https://tmpnotes.com). The website has an encryption feature that can encrypt your note in the browser *before* it is sent to the server. This ensures that even if someone did view the database, the data would not be readable. The codebase is also small, and relatively easy to audit.

If you don't want to send your information to [tmpnotes.com](https://tmpnotes.com), we don't blame you! We want to make this project easy to run yourself if that is desired. We will host [tmpnotes.com](https://tmpnotes.com) as long as it does not become frequently abused or prohibitively expensive for us to do so.

# Run it yourself
We understand that secrets are sensitive and people may not want to use a publicly hosted instance ([tmpnotes.com](https://tmpnotes.com)). The following is information you can use to run a TMPNOTES system yourself! We strongly suggest that any implementation you run yourself is fronted by a reverse proxy (such as NGINX) with TLS.

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


