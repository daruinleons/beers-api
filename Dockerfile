FROM golang:alpine AS build

ENV GOPROXY=https://proxy.golang.org

WORKDIR /go/src/beers-api
COPY . .

RUN apk add build-base
RUN GOOS=linux go build -o /go/bin/beers-api src/main.go

FROM alpine
COPY --from=build /go/bin/beers-api /go/bin/beers-api
ENTRYPOINT ["/go/bin/beers-api"]