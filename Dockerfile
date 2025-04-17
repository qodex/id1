FROM golang:latest AS build

WORKDIR /go/src/app
COPY *.go .
COPY go.* .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -ldflags="-X main.version=$(date +%Y%m%d)" -o /go/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]