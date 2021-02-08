FROM  golang:1.14-buster as builder

WORKDIR /tmp/build
COPY . ./

RUN go mod tidy 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o main src/app/backend/main.go

FROM scratch
COPY --from=builder /tmp/build /
CMD ["/main"]