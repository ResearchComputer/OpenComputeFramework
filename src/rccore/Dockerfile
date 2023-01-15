FROM golang:alpine

ENV GIN_MODE=release
ENV PORT=8092
COPY . /app/
WORKDIR /app

RUN apk add --no-cache git make bash && make build

EXPOSE $PORT

ENTRYPOINT ["/app/build/rccore", "start", "--config", "/app/.rccore/cfg.yaml"]