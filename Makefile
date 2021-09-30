STACK_NAME=lambdaapi

build:
	cd cmd/api && GOARCH=amd64 GOOS=linux go build -o ../../dist/api

format:
	gofmt -s -w .

validate:
	sam validate -t api.yaml

run:
	go run cmd/api/main.go

deploy:
	sam deploy \
		-t api.yaml \
		--stack-name=$(STACK_NAME) \
		--capabilities=CAPABILITY_IAM \
		--resolve-s3

deploy-local:
	samlocal validate -t api.yaml --region=us-east-1
	samlocal deploy \
		-t api.yaml \
		--stack-name=$(STACK_NAME) \
		--capabilities=CAPABILITY_IAM \
		--region=us-east-1 \
		--resolve-s3

get-apigw-local-id:
	$(eval API_ID=$(shell awslocal apigateway get-rest-apis | jq ".items[] | .id"))
	@echo $(API_ID)