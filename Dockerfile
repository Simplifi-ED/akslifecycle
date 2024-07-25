FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.5 AS build-stage

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on

WORKDIR /akslifecycle
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=.,target=.,rw \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /tmp/akslifecycle -ldflags="-s -w"


FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.19 AS build-release-stage

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates \
    && apk cache clean

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    akslifecycle

USER akslifecycle

WORKDIR /app

COPY --from=build-stage /tmp/akslifecycle ./

ENTRYPOINT ["/app/akslifecycle"]
