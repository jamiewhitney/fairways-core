FROM golang:latest as builder
ARG SERVICE
WORKDIR /app
COPY . .
COPY ./cmd/${SERVICE} .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o main .


FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 3000
CMD ["./main"]