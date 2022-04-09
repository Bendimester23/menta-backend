FROM golang:latest AS build

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .

FROM centurylink/ca-certs

WORKDIR /app

#TODO: remove comment if needed
#RUN mkdir ./static
#COPY ./static ./static

ENV PORT=8080
ENV MODE=prod
ENV EMAIL_HOST=172.19.0.2

COPY --from=build /go/src/app/app .

CMD ["./app"]