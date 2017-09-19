FROM golang:1.9

WORKDIR /go/src/github.com/sjeandeaux/nexus-cli
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ENTRYPOINT ["go-wrapper" , "run"]