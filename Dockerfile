FROM golang:1.22 AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /todolist ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /todolist /todolist
EXPOSE 8080
USER nonroot:nonroot
CMD ["/todolist"]