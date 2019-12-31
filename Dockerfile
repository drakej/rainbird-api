FROM golang:1.13.4-alpine as build

LABEL maintainer="Jonathan Drake <drakej@drakej.com>"

WORKDIR /go/src/github.com/drakej/rainbird-api

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine:latest

WORKDIR /opt/rainbird/

COPY --from=build /go/src/github.com/drakej/rainbird-api/config.toml .
COPY --from=build /go/src/github.com/drakej/rainbird-api/sipCommands.json .
COPY --from=build /go/src/github.com/drakej/rainbird-api/rainbird-api .

EXPOSE 8080

CMD ["./rainbird-api"]