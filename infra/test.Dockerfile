# Build the Go Binary.
FROM golang:1.18 as build_test
ENV CGO_ENABLED 0
ARG BUILD_REF

# Create the challenge directory and the copy the module files first and then
# download the dependencies.
RUN mkdir /challenge
COPY go.* /challenge/
WORKDIR /challenge
RUN go mod download

# Copy the source code into the container.
COPY . /challenge

# Build the challenge binary. We are doing this last since this will be different
# every time we run through this process.
WORKDIR /challenge/cmd/challenge
CMD CGO_ENABLED=1 go test -count=1 -v ./...
