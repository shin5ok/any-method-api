FROM golang:1.19 AS builder

WORKDIR /app
COPY . .
RUN GGO_ENABLED=0 GOOS=linux go build -o main

FROM golang:1.19 AS runner
# FROM busybox AS runner
COPY --from=builder /app/main /main
RUN rm -f /dev/random && ln -s /dev/urandom /dev/random
USER nobody
ENV GIN_MODE release
CMD ["/main"]

