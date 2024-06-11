FROM golang:1.22.0 AS build-stage
WORKDIR /app
COPY load-shedding/cmd/main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /service main.go
FROM gcr.io/distroless/base AS build-release-stage
WORKDIR /
COPY --from=build-stage /service /service
EXPOSE 8002
ENTRYPOINT [ "/service" ]