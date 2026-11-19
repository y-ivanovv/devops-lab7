FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/todo-service ./cmd/todo-service
COPY internal ./internal
COPY proto ./proto
RUN CGO_ENABLED=0 go build -o /todo-service ./cmd/todo-service

FROM scratch
COPY --from=build /todo-service /todo-service
USER 65532:65532
EXPOSE 8080 9090
ENTRYPOINT ["/todo-service"]
