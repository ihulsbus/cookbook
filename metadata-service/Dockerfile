############################
# STEP 1 build executable binary
############################
FROM golang AS builder

# Ensure ca-certficates are up to date
RUN update-ca-certificates

WORKDIR /build

# use modules
COPY . .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify


# Build the binary.
RUN CGO_ENABLED=0 go build cmd/metadata-service/main.go

############################
# STEP 2 build a small image
############################
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER 1000
COPY --from=builder /build/main /go/bin/main

ENTRYPOINT [ "/go/bin/main" ]
