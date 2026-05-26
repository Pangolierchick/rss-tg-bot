FROM golang:1.26 AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum ./
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/cc-debian12:nonroot AS runtime
COPY --from=build /app/bin/rss-tg-bot /rss-tg-bot

USER nonroot:nonroot

CMD ["/rss-tg-bot"]
