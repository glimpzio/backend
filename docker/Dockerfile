FROM golang:1.20.5-alpine3.17 AS build

WORKDIR /usr/build

COPY src/go.sum src/go.mod ./

RUN ["go", "mod", "download"]

COPY src/ ./

RUN ["go", "build", "-o", "main", "main.go"]

FROM alpine:3.17 AS runtime

ENV GIN_MODE=release

WORKDIR /usr/app

COPY --from=build /usr/build/main ./

CMD ["./main"]