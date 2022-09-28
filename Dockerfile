FROM golang:1.16 AS builder

WORKDIR /app
COPY . .
RUN go build -o main
# RUN GGO_ENABLED=0 GOOS=linux go build -o main

FROM busybox AS runner
COPY --from=builder /app/main /main
RUN rm -f /dev/random && ln -s /dev/urandom /dev/random
USER nobody
CMD ["/main"]

