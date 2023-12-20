FROM golang:1.21.5-alpine AS build-stage

WORKDIR /akslifecycle

ARG CGO_ENABLED=0 

COPY go.mod go.sum ./
COPY internal /internal
COPY cmd /cmd
COPY utils /utils
COPY main.go /

RUN go mod download && go build -o /akslifecycle

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage  akslifecycle /

ENTRYPOINT ["/akslifecycle"]
