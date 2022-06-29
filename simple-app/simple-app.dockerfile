# base image
FROM golang:1.18-alpine as builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o simpleApp
RUN chmod +x /app/simpleApp

# small image with just executable
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/simpleApp /app
CMD ["/app/simpleApp"]