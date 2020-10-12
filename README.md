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

## Features

### For models

    detailed-pagination: true - добавляет count в response метода list

    id-from-profile: true - позволяет совершать методы для текущего пользователя

    return-when-edit: true - вобавляет модель в response метода edit

    tags:
    - Tag - добавит тег Tag  во все эндпоинты модели

### For columns

    sort-on: true - позволяет делать сортировку по данной колонке в методе list

### For methods

    edit(column2, column3) - в таком виде метод edit будет изменять только column2 и column3, а также находится в эндпоинте /model/editColumn2Column3

    list(column1, column3*, model1*(column1, model1(column1, column2))) - в таком виде метод list будет возвращать только column2, column3 и model1..., то есть только рекурсивно указанные поля и фильтрация может проиваодиться по column3 и по id model1, а также находится в эндпоинте /model/editColumn2Column3Model1
