FROM golang:1.11.9-alpine3.8 AS build_base
RUN apk add bash make git curl unzip rsync libc6-compat gcc musl-dev
WORKDIR /go/src/github.com/spacemeshos/connector

# Force the go compiler to use modules
ENV GO111MODULE=on

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

COPY scripts/* scripts/

RUN ./scripts/setup_env.sh

FROM build_base AS server_builder
COPY . .

RUN ./scripts/genproto.sh
RUN cd main && go build .

FROM alpine AS connector 

# Finally we copy the statically compiled Go binary.
COPY --from=server_builder /go/src/github.com/spacemeshos/connector/main/main /bin/main

ENTRYPOINT ["/bin/main"]

