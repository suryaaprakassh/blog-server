FROM golang:1.24-alpine AS build
WORKDIR /app

COPY go.mod go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./main .

FROM  alpine:3.14 AS release
WORKDIR /
COPY --from=build /app/main /main
COPY --from=build /app/static /static
COPY --from=build /app/posts /posts
EXPOSE 8000

ENTRYPOINT ["/main"]
