FROM golang:alpine as builder
WORKDIR /urlshortner
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o urlshortner cmd/*.go

FROM scratch
COPY --from=builder /urlshortner/urlshortner .
CMD ["./urlshortner", "serve"]