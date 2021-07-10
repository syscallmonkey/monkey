FROM golang:1.16-alpine as builder

# Install build deps

RUN apk add --update git make bash

# Get dependencies

WORKDIR /w
COPY go.mod go.sum /w/
RUN go mod download

# Build

COPY . ./
RUN make bin/monkey

# Build the asset container, copy over the binary

FROM scratch
COPY --from=builder /w/bin/monkey /monkey
ENTRYPOINT ["/monkey"]
