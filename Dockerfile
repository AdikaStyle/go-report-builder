FROM golang:1.13-alpine AS builder

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
  -installsuffix 'static' \
  -o /app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app .

RUN mkdir /templates

ENV TEMPLATES_PATH="/templates"
ENV RENDER_TIMEOUT="3s"
ENV SERVER_HOST="localhost"
ENV SERVER_PORT="8080"
ENV GIN_MODE="release"

CMD ["./app"]