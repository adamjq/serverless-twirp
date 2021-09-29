S3_BUCKET=artifacts-aws-sam
STACK_NAME=lambdaapi

build:
	cd cmd/api && GOARCH=amd64 GOOS=linux go build -o ../../dist/api

format:
	gofmt -s -w .

run:
	go run cmd/api/main.go

deploy:
	sam deploy \
		-t api.yaml \
		--stack-name=$(STACK_NAME) \
		--s3-bucket=$(S3_BUCKET) \
		--capabilities=CAPABILITY_IAM