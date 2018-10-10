FROM golang:1.11

WORKDIR /app

ADD . .
RUN go get ./...

RUN mkdir /out && CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o /out/kms-operator ./cmd/kms-operator/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=0 /out/kms-operator /app/kms-operator
ENTRYPOINT ["/app/kms-operator"]