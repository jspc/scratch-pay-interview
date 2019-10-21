FROM golang:1.13.0 as build

WORKDIR /app
ADD . .

RUN go get -v && \
    CGO_ENABLED=0 go build

FROM scratch
COPY --from=build /usr/share/ca-certificates /usr/share/ca-certificates

ENV LISTEN_ADDR=0.0.0.0:8080

CMD ["/app"]
COPY --from=build /app/app /
