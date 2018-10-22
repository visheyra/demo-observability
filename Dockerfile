FROM golang:1.11.1-stretch as build
WORKDIR $GOPATH/src/github.com/visheyra/demo-observability/
ADD . .
RUN go install github.com/visheyra/demo-observability

FROM gcr.io/distroless/base
COPY --from=build /go/bin/demo-observability /app
ENTRYPOINT ["/app"]
