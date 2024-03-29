FROM golang:1.17-alpine as build
WORKDIR /app
ADD . .
ENV GOPATH /go
ENV CGO_ENABLED=0

ARG DISABLE_TESTS
RUN if [[ "$DISABLE_TESTS" = "true" ]] ; then echo Skipping Tests ; else go test ./...; fi

ARG VERSION
ARG GITSHA
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'tmpnotes/internal/version.version=${VERSION}' -X 'tmpnotes/internal/version.gitSHA=${GITSHA}'"

FROM alpine:latest
COPY --from=build /app/tmpnotes /app/
WORKDIR /app
RUN chown 65534:65534 tmpnotes
USER 65534:65534
ADD static/ ./static
ADD templates/ ./templates
ENTRYPOINT [ "./tmpnotes" ]
