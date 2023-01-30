FROM golang:latest

RUN apt update && apt -y upgrade && apt install sudo

WORKDIR /usr/src/mt-gen
COPY ./ ./

RUN GO111MODULE=on go get github.com/golang/mock/mockgen && bash install

ENTRYPOINT ["mt-gen"]