CDK_DIR=_cdk

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

run:
	go run cmd/api/main.go

########## CDK CI ##########

cdk-install:
	cd $(CDK_DIR) && npm ci

cdk-lint:
	cd $(CDK_DIR) && npm run lint && npm run prettier

cdk-test:
	cd $(CDK_DIR) && npm run test

cdk-build:
	cd $(CDK_DIR) && npm run build

########## AWS ##########

AWS_VAULT := aws-vault exec $(AWS_PROFILE) --

cdk-bootstrap:
	cd $(CDK_DIR) && $(AWS_VAULT) npm run cdk bootstrap

cdk-deploy: build
	cd $(CDK_DIR) && $(AWS_VAULT) npm run cdk synth && $(AWS_VAULT) npm run cdk deploy

cdk-teardown:
	cd $(CDK_DIR) && $(AWS_VAULT) npm run cdk destroy