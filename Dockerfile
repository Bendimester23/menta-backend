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
ENV EMAIL_HOST=zeus.bendi.cf
ENV EMAIL_USER=apikey
ENV EMAIL_PASS=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
ENV EMAIL_PORT=1025

COPY --from=build /go/src/app/app .

EXPOSE ${PORT}

CMD ["./app"]