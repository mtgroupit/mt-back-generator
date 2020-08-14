# mt-back-generator

## Preinstallation

For using this app install:

1. [Golang](https://golang.org/doc/install)

2. [Docker](https://docs.docker.com/get-docker/) and execute [Post-installation steps](https://docs.docker.com/engine/install/linux-postinstall/)

3. [Gomock](https://github.com/golang/mock) using command:

    ```bash
    GO111MODULE=on go get github.com/golang/mock/mockgen
    ```

## Installation

```bash
git clone https://github.com/mtgroupit/mt-back-generator.git
cd mt-back-generator
make
```

## Using

```bash
make gen [TARGET_DIRECTORY=~/Desktop] [CONFIG=./config.yaml]
```

Open TARGET_DIRECTORY/service.
