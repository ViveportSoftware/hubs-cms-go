############################
# STEP 1 BUILD EXECUTABLE BINARY
############################
FROM golang:1.17-alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
# create appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
# fetch dependencies.
RUN go mod download
RUN go mod verify
RUN go get github.com/swaggo/swag/cmd/swag@v1.7.3
RUN swag init --parseDependency --parseInternal
# build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/hubs-cms-go
############################
# STEP 2 BUILD A SMALLER IMAGE
############################
FROM scratch
# pass FULL_VERSION to environment variable
ARG FULL_VERSION
ENV FULL_VERSION ${FULL_VERSION:-unknown}
# import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# copy executable
COPY --from=builder /go/bin/hubs-cms-go /go/bin/hubs-cms-go
# use an unprivileged user.
USER appuser:appuser
# expose port
EXPOSE 9999
# run binary.
ENTRYPOINT ["/go/bin/hubs-cms-go"]
