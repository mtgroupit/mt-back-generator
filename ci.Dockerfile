FROM golang:1-buster


RUN apt-get update -yqq && \
    apt-get install -yqq sudo git build-essential && \
    apt-get clean && \
    apt-get autoremove -yqq && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY . /app/
WORKDIR /app

RUN go get github.com/golang/mock/mockgen && \
    bash ./install

RUN mt-gen -dir=./generated/ -config=./samples/config_mini_demo.yaml

WORKDIR /app/extauthapi/
RUN go test ./...

# TODO
#    ./test && \
#    ./testall