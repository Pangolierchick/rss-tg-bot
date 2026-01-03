FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/static-debian13 AS runtime
COPY --from=build /app/bin/rss-tg-bot /rss-tg-bot

CMD ["/rss-tg-bot"]
