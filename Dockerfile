FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.21.5 AS build-stage

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on

WORKDIR /akslifecycle
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ./akslifecycle

FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch AS build-release-stage

WORKDIR /app
COPY --from=build-stage /akslifecycle/akslifecycle ./

ENTRYPOINT ["/app/akslifecycle"]
