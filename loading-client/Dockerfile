FROM golang:1.18 AS builder
RUN apt update && apt -y install git
RUN mkdir -p /app/src
RUN git clone https://github.com/rakyll/hey.git /app/src/hey
RUN cd /app/src/hey; go build -o hey

FROM debian:11-slim AS runner
COPY --from=builder /app/src/hey/hey /usr/bin/
COPY run.sh /
CMD ["/run.sh", "2"]
