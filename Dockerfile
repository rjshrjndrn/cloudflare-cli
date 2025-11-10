FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY cfcli /usr/local/bin/cfcli

ENTRYPOINT ["/usr/local/bin/cfcli"]
