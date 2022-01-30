FROM golang:1.17-buster as base

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

COPY . .

FROM base as builder

RUN make build

FROM alpine:latest

COPY --from=builder /app/bin/ /app/bin/
WORKDIR /app
EXPOSE 8080
CMD ["./bin/pet-me-api"]
