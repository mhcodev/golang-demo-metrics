FROM golang:1.23.5-alpine3.21 AS build

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o /app/main .

FROM alpine:3.21

WORKDIR /app

COPY --from=build /app/main /app/main
CMD [ "/app/main" ]