FROM golang:1.21 as builder

ENV ENVIRONMENT=prod

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/server cmd/server/main.go

FROM alpine:3.14

ENV ENVIRONMENT=prod

ENV APP_PORT=${APP_PORT}
ENV MASTER_KEY=${MASTER_KEY}
ENV UPSTASH_USERNAME=none
ENV UPSTASH_PASSWORD=none
ENV UPSTASH_URL=none
ENV EMAIL_LEADS_SHEET_ID=${EMAIL_LEADS_SHEET_ID}

WORKDIR /opt/app

COPY --from=builder /src/server /opt/app/server
COPY --from=builder /src/config/app.env /opt/app/config/app.env
COPY --from=builder /src/config/ijnrge.env.enc /opt/app/config/ijnrge.env.enc

ENTRYPOINT ["/opt/app/server"]
