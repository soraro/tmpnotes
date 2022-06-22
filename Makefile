BINARY_NAME := tmpnotes
GITSHA := $(shell git rev-parse HEAD)
GITTAG := $(shell git describe --tags)

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	@go fmt
	@if test -z $(gofmt -l .); then \
		echo "All golang files formatted correctly 👍️"; \
	else \
		echo "❗️ Golang formatting issues:"; gofmt -l .; \
	fi

.PHONY: build
build:
	go build -o ${BINARY_NAME} -ldflags "-s -w -X 'tmpnotes/internal/version.version=${GITTAG}' -X 'tmpnotes/internal/version.gitSHA=${GITSHA}'"

.PHONY: local-env
local-env:
	docker-compose up

# just run a redis container for local development
.PHONY: dev
dev:
	docker run -d --rm -p 6379:6379 --name tmpnotes_local redis

.PHONY: dev-down
dev-down:
	docker stop tmpnotes_local

# Build a container image locally for development
.PHONY: container-dev
container-dev:
	docker-compose -f docker-compose-dev.yaml build --build-arg VERSION=${GITTAG} --build-arg GITSHA=${GITSHA}
	docker-compose -f docker-compose-dev.yaml up

.PHONY: container-dev-down
container-dev-down:
	docker-compose -f docker-compose-dev.yaml down

.PHONY: deploy-heroku
deploy-heroku:
	heroku container:login
	heroku container:push web --arg VERSION=${GITTAG},GITSHA=${GITSHA} -a tmpnotes
	heroku container:release web -a tmpnotes