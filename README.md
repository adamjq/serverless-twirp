# serverless-rpc
Serverless Golang RPC API with MVC structure

## Dependencies:
- AWS CLI
- SAM CLI
- [AWS SAM Local CLI](https://github.com/localstack/aws-sam-cli-local)
- Docker
- Go 1.17

## Development

### Unit tests

```bash
make test
```

### Manual testing

```bash
docker-compose up
```

In a separate terminal window run:
```
make deploy-local

awslocal s3 ls
awslocal lambda list-functions

make get-apigw-local-id
curl -X POST http://localhost:4566/restapis/{API_ID}/dev/hello -d 'Hello, world!'
```

Localstack services can be inspected at http://localhost:4566/health.

## Deployment

```bash
make deploy
```