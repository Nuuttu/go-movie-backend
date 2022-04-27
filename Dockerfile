FROM golang:1.18

WORKDIR /app
COPY ./endpoints ./app
COPY ./structs ./app
COPY ./utils ./app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN go build -v -o /go-movie-docker

EXPOSE 10000

CMD ["/go-movie-docker"]