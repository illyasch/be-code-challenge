# Build the Go Binary.
FROM golang:1.18 as build_challenge
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
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# Run the Go Binary in Alpine.
FROM alpine:3.15
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S challenge && \
    adduser -u 1000 -h /challenge -G challenge -S challenge
COPY --from=build_challenge --chown=challenge:challenge /challenge/cmd/challenge/challenge /challenge/challenge
WORKDIR /challenge
USER challenge
CMD ["./challenge"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="challenge" \
      org.opencontainers.image.authors="Ilya Scheblanov <ilya.scheblanov@gmail.com>" \
      org.opencontainers.image.source="https://github.com/illyasch/be-code-challenge/cmd/challenge" \
      org.opencontainers.image.revision="${BUILD_REF}"
