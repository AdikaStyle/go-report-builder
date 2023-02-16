FROM ubuntu:22.04

WORKDIR /go/

COPY . .

RUN apt-get update -y && apt-get install -y golang nodejs npm 

RUN cd /go/ui && npm install && npx playwright install-deps chromium

RUN cd ui && npm run build

RUN go build -o /bin/greypot-server cmd/greypot-server/*.go

CMD [ "/bin/greypot-server" ]