FROM golang:1.18

WORKDIR /app

COPY go.mod go.sum ./

COPY ./endpoints ./endpoints
COPY ./structs ./structs
COPY ./utils ./utils

COPY ./Medialists.xlsx ./

RUN go mod download
RUN go mod verify

COPY *.go ./
RUN go build -v -o /go-movie-docker

EXPOSE 10000

CMD ["/go-movie-docker"]