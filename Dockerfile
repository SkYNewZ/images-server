FROM golang:1.16
WORKDIR /go/src/github.com/SkYNewZ/images-server

COPY go.* ./
RUN go mod download

COPY . .
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    export COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null) && \
    export VERSION=$(git describe --tags --exact-match 2>/dev/null || git describe --tags 2>/dev/null || echo "v0.0.0-${COMMIT_HASH}") && \
    go build -ldflags "-X 'github.com/SkYNewZ/images-server/internal.buildNumber=${VERSION}'" -o /images-server .


FROM scratch

ARG BUILD_DATE
ARG VCS_REF

ENV PORT 8080

LABEL maintainer="Quentin Lemaire <quentin@lemairepro.fr>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.build-date=$BUILD_DATE
LABEL org.label-schema.name="skynewz/images-server"
LABEL org.label-schema.description="REST API to manage pictures stored in a S3-compatible backend"
LABEL org.label-schema.vcs-url="https://github.com/SkYNewZ/images-server"
LABEL org.label-schema.vcs-ref=$VCS_REF

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=0 /images-server /images-server
EXPOSE ${PORT}
ENTRYPOINT ["/images-server"]