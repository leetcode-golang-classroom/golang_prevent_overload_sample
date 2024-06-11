FROM golang:1.22.0 AS build-stage
WORKDIR /app
COPY rate-limiter/cmd/main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /service main.go
FROM gcr.io/distroless/base AS build-release-stage
WORKDIR /
COPY --from=build-stage /service /service
EXPOSE 8001
ENTRYPOINT [ "/service" ]