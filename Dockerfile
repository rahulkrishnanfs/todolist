FROM golang:1.25 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /todolist ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /todolist /todolist
COPY config/properties.toml ./config/properties.toml 
EXPOSE 8080
USER nonroot:nonroot
CMD ["/todolist"]