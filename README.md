# serverless-twirp
Serverless Golang RPC API with MVC structure

## Why Twirp?



## Why Lambda?

## Dependencies:
- AWS CLI
- AWS CDK Local CLI](https://github.com/localstack/aws-cdk-local)
- Docker
- direnv
- Go 1.17
- [buf](https://docs.buf.build/installation/)

## Development

### Unit tests

```bash
make ci
```

### Localstack

```bash
docker-compose up
```

In a separate terminal window run:
```
make local-cdk-bootstrap
make local-cdk-deploy

# test function is deploy
awslocal lambda list-functions
```

Localstack services can be inspected at http://localhost:4566/health.

## Deployment

The app uses [aws-vault](https://github.com/99designs/aws-vault) to deploy to AWS environments.

```
direnv allow .

# bootstrap AWS environment if not already done
make cdk-bootstrap

make cdk-deploy
```

## References

- [Alex DeBrie DynamoDB blog](https://www.alexdebrie.com/posts/dynamodb-no-bad-queries/)
- [DynamoDB AWS Go V2 cheatsheet](https://dynobase.dev/dynamodb-golang-query-examples/)
