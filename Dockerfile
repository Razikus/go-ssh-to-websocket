
FROM golang:1.22.2-alpine3.19 AS builder
COPY src /src
WORKDIR /src
RUN go build -o /sshtows

FROM scratch
COPY --from=builder /sshtows /sshtows
COPY --from=builder /src/index.html /index.html
COPY --from=builder /src/xterm.js /xterm.js
COPY --from=builder /src/xterm.css /xterm.css

EXPOSE 8280
ENTRYPOINT [ "/sshtows" ]