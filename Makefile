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

########## AWS ##########

AWS_VAULT := aws-vault exec $(AWS_PROFILE) --

cdk-bootstrap:
	cd _cdk && $(AWS_VAULT) npm run cdk bootstrap

cdk-deploy: build
	cd _cdk && $(AWS_VAULT) npm run cdk synth && $(AWS_VAULT) npm run cdk deploy

cdk-teardown:
	cd _cdk && $(AWS_VAULT) npm run cdk destroy