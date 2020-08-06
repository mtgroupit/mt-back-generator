TARGET_DIRECTORY=${1:-~/Desktop}
CONFIG=${2:-./config.yaml}

echo "Generator started"
./bin/generator --config $CONFIG --dir $TARGET_DIRECTORY
echo "Generator stopped"

echo "fmt started"
cd $TARGET_DIRECTORY/service
go fmt ./...
echo "fmt stopped"

echo "Swagger generate started"
alias swagger="docker run --rm -it -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger"
swagger generate server -t $TARGET_DIRECTORY/service/internal/api/restapi -f $TARGET_DIRECTORY/service/swagger.yaml --exclude-main
swagger generate client -t $TARGET_DIRECTORY/service/internal/api/restapi -f $TARGET_DIRECTORY/service/swagger.yaml
echo "Swagger generate  stopped"

echo "Mock generate started"
cd $TARGET_DIRECTORY/service/internal/app
mockgen -source=app.go -destination=testing.generated.go -package=app
echo "Mock generate stopped"