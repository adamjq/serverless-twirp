generate:
	protoc \
		--go_out=./internal/userpb \
	  --twirp_out=./internal/userpb \
		./proto/user/v1/user.proto
	@go generate ./...

build: generate
	cd cmd/api && GOARCH=amd64 GOOS=linux go build -o ../../dist/api

test: generate
	go test ./...

lint:
	buf lint

format:
	gofmt -s -w .

ci: lint test

run:
	go run cmd/api/main.go

########## AWS ##########

AWS_VAULT := aws-vault exec $(AWS_PROFILE) --

cdk-bootstrap:
	cd _cdk && $(AWS_VAULT) npm run cdk bootstrap

cdk-deploy: build
	cd _cdk && $(AWS_VAULT) npm run cdk synth && $(AWS_VAULT) npm run cdk deploy

########## LOCALSTACK ##########

get_local_api_id = $(shell awslocal apigateway get-rest-apis | jq ".items[] | .id"| tr -d '"')

local-cdk-bootstrap:
	cd _cdk && npm run cdklocal bootstrap

local-cdk-deploy: build
	cd _cdk && npm run cdklocal synth && npm run cdklocal deploy

get-apigw-local-id:
	$(eval APIGW=http://$(get_api_id).execute-api.localhost.localstack.cloud:4566/prod/)
	@echo $(APIGW)

call-api-get:
	curl http://$(get_api_id).execute-api.localhost.localstack.cloud:4566/prod/twirp/proto.user.v1.UserService/GetUser \
		-X POST \
		-H 'Content-Type: application/json' \
		-d '@testdata/GetUserRequest.json'

call-api-store:
	curl http://$(get_api_id).execute-api.localhost.localstack.cloud:4566/prod/twirp/proto.user.v1.UserService/StoreUser \
		-X POST \
		-H 'Content-Type: application/json' \
		-d '@testdata/StoreUserRequest.json'
