FROM golang:1.20.4-alpine3.18 AS builder

RUN go version
COPY ./ /github.com/andredubov/todo-backend
WORKDIR /github.com/andredubov/todo-backend

RUN go mod download && go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./.bin/app ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root

COPY --from=builder /github.com/andredubov/todo-backend/.bin/app .
ADD configs ./configs

ARG TAG
ENV TAG ${TAG}

CMD [ "./app"]