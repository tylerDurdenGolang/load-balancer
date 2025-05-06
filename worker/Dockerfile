# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mathcruncher .

FROM scratch
COPY --from=builder /app/mathcruncher /mathcruncher
EXPOSE 8080
ENTRYPOINT ["/mathcruncher"]
