package mocks

//go:generate env GOBIN=$PWD/bin GO111MODULE=on go install github.com/golang/mock/mockgen

// Store mocks
//go:generate $PWD/bin/mockgen -destination storemock/storemock.go -package storemock github.com/adamjq/serverless-twirp/internal/stores Users
