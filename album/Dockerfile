FROM golang:1.27 AS build

WORKDIR /go/src/app
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download && go mod verify

COPY *.go ./

RUN go vet -v
RUN go test -v

RUN go build -o /go/bin/app

FROM gcr.io/distroless/static:nonroot

COPY --from=build /go/bin/app /app
CMD ["/app"]
