# mt-back-generator

## Preinstallation

For using this app install:

1. [Golang](https://golang.org/doc/install)

2. [Docker](https://docs.docker.com/get-docker/) and execute [Post-installation steps](https://docs.docker.com/engine/install/linux-postinstall/)

3. [Gomock](https://github.com/golang/mock) using command:

    ```bash
    GO111MODULE=on go get github.com/golang/mock/mockgen
    ```

4. [Connecting to GitHub with SSH](https://docs.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh)

## Installation

```bash
git clone git@github.com:mtgroupit/mt-back-generator.git
cd mt-back-generator
bash install
```

## Using

1. Write the configuration
    - [Documentation](https://www.notion.so/mtgroupit/4a109d202c0443e6b222d33dff3a2e4e)
    - [Examples](./samples)

2. Generate the project from configuration

    ```bash
    mt-gen -config=./samples/config_mini_demo.yaml
    ```

    Options (flags)
    - config - the path to the configuration (Default: "./config.yaml")
    - dir - the target dir for the generated project (Default: "./generated/")


3. Go to the folder with the generated project and run the build script

    ```bash
    bash build
    ```

That's it, the project is ready to launch and use!
