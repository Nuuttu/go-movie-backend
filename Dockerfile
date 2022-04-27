FROM golang:1.18

# RUN mkdir /go/app

# COPY ./endpoints ./app/endpoints
# COPY ./structs ./app/structs
# COPY ./utils ./app/utils

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change

COPY go.mod go.sum ./

COPY ./endpoints ./endpoints
COPY ./structs ./structs
COPY ./utils ./utils

COPY ./Medialists.xlsx ./

RUN go mod download
RUN go mod verify
# RUN go get

COPY *.go ./
RUN go build -v -o /go-movie-docker

EXPOSE 10000

CMD ["/go-movie-docker"]