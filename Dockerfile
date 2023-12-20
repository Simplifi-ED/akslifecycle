FROM golang:1.21.5 AS build-stage

WORKDIR /akslifecycle
COPY go.mod go.sum /
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./akslifecycle

FROM scratch AS build-release-stage

WORKDIR /app
COPY --from=build-stage /akslifecycle/akslifecycle /

ENTRYPOINT ["/akslifecycle"]
