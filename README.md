# serverless-rpc
Serverless Golang RPC API with MVC structure

## Dependencies:
- AWS CLI
- SAM CLI
- Docker
- Go 1.17

## Deployment

An S3 bucket first needs to be created in AWS to store lambda artifacts. Update the `S3_BUCKET` variable in the Makefile
with the bucket name.

```bash
make deploy
```