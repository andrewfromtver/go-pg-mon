ARG ALPINE_VER="3.20.3"
ARG GO_VER="1.23.0"

# Stage 1: Build Go app
FROM golang:${GO_VER}-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o pg_query_api .

# Stage 2: Create image
FROM alpine:${ALPINE_VER}

RUN addgroup -S api_group && adduser -S api_user -G api_group
WORKDIR /home/api_user/
COPY --from=builder /app/pg_query_api .
RUN chown api_user:api_group /home/api_user/pg_query_api
USER api_user

EXPOSE 8080

CMD ["./pg_query_api"]

