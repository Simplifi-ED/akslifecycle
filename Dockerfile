FROM golang:1.21.5-alpine AS build-stage

WORKDIR /akslifecycle

COPY go.mod go.sum /
RUN go mod download
COPY . /

RUN CGO_ENABLED=0 go build -o /akslifecycle

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app

COPY --from=build-stage  akslifecycle /

ENTRYPOINT ["/akslifecycle"]
