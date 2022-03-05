#!/bin/bash

set -eu

aws dynamodb --endpoint-url ${DYNAMO_ENDPOINT} create-table \
    --table-name FakeTable \
    --attribute-definitions AttributeName=PK,AttributeType=S AttributeName=SK,AttributeType=S \
    --key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    || true # don't fail if table already exists
