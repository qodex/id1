FROM golang:latest AS build

WORKDIR /go/src/app
COPY main.go_ ./main.go

RUN go mod init id1.au/id1
RUN go get github.com/joho/godotenv
RUN go get github.com/qodex/id1
RUN go get -u all

RUN CGO_ENABLED=0 go build -ldflags="-X main.version=$(date +%Y%m%d)" -o /go/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]