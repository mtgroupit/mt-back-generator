// Package extauthapi defines external Swagger API for extauth.
package extauthapi

//go:generate rm -rf model restapi client
//go:generate sh -c "cd '$MIRROR_IN_GOPATH' && swagger generate server --template-dir=templates --api-package op --model-package model --exclude-main --principal extauthapi.Profile --strict --existing-models=extauth/internal/extapi/mtmb-extauthapi"
//go:generate sh -c "cd '$MIRROR_IN_GOPATH' && swagger generate client --template-dir=templates --api-package op --model-package model"
// Cut repo directory name from this file, to ensure it won't change when
// generated inside CircleCI in default working_directory (~/project).
//go:generate sh -c "perl -i -pe 's,\\.\\./\\.\\./\\w\\S*,..,' restapi/configure_*.go"
