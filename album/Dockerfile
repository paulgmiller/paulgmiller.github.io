FROM golang:1.23 as build

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download && go mod verify

COPY *.go ./

RUN go vet -v
RUN go test -v

RUN go build -o /go/bin/app

FROM gcr.io/distroless/base

COPY --from=build /go/bin/app /
CMD ["/app"]
