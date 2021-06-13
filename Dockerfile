FROM golang:1.13-alpine AS builder

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
  -installsuffix 'static' \
  -o /app .

FROM alpine:latest

RUN apk update && apk add --no-cache nmap && \
    echo @edge http://nl.alpinelinux.org/alpine/edge/community >> /etc/apk/repositories && \
    echo @edge http://nl.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories && \
    apk update && \
    apk add --no-cache \
      ca-certificates \
      chromium \
      harfbuzz \
      "freetype>2.8" \
      ttf-freefont \
      nss

COPY --from=builder /app .

RUN mkdir /templates

ENV TEMPLATES_PATH="/templates"
ENV RENDER_TIMEOUT="5s"
ENV SERVER_HOST="localhost"
ENV SERVER_PORT="8080"
ENV GIN_MODE="release"

ADD examples /templates

CMD ["./app"]