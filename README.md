# mt-back-generator

## Preinstallation

For using this app install:

1. [Golang](https://golang.org/doc/install)

1. [Docker](https://docs.docker.com/get-docker/) and execute [Post-installation steps](https://docs.docker.com/engine/install/linux-postinstall/)

2. [Gomock](https://github.com/golang/mock) using command:
    ```
    GO111MODULE=on go get github.com/golang/mock/mockgen
    ```

## Installation

```
git clone https://github.com/mtgroupit/mt-back-generator.git
cd mt-back-generator
make
```

## Using

```
make gen [TARGET_DIRECTORY=~/Desktop] [CONFIG=./config.yaml]
```
