FROM golang:1.11.2-alpine AS builder


WORKDIR /go/src/github.com/mteodor/edgex-app
COPY . .

RUN apk update \
    && apk add make \
    && mv build/edgex-app /exe

FROM scratch
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
