# Swagger-описание API сервиса extauth и сгенерированный из него код
[![CircleCI](https://circleci.com/gh/mtgroupit/mtmb-extauthapi.svg?style=svg&circle-token=28c95158bcebf28b20f60dcee817b342324567dd)](https://circleci.com/gh/mtgroupit/mtmb-extauthapi)


## Setup
### Go tools
Install tools required to build/test project.
```
go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.18.0
git clone https://github.com/Djarvur/go-swagger.git &&
    cd go-swagger/ && go install ./cmd/swagger/ && cd .. && rm -rf go-swagger/
```


## Testing
- `go test ./...` - test project, fast
- `./test` - carefully test project
- `./testall` - carefully test project including integration tests in
  exactly same way as CI, slow
