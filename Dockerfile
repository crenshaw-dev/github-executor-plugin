FROM golang:1.19-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w" -o plugin cmd/github-plugin/main.go

FROM gcr.io/distroless/static

COPY --from=build /app/plugin /

CMD [ "/plugin" ]
