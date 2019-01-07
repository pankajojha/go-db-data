FROM golang:alpine

#ADD ./ /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app

ENV PORT=3001
ENV GO111MODULE=on
#RUN go mod init
RUN go get -d -v
#RUN go get github.com/go-sql-driver/mysql
RUN go mod download

RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app .

CMD ["go", "run", "main.go"]