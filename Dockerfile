FROM golang:1.25.0-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tilokit .

FROM alpine:latest

RUN apk --no-cache add ca-certificates git nodejs npm
RUN addgroup -g 1001 -S tilokit && \
	adduser -S tilokit -u 1001 -G tilokit

WORKDIR /workspace

COPY --from=builder /app/tilokit /usr/local/bin/tilokit

RUN chmod +x /usr/local/bin/tilokit
RUN chown -R tilokit:tilokit /workspace

USER tilokit

ENTRYPOINT ["tilokit"]
CMD ["--help"]
