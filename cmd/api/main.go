package main

import (
	"context"
	"log"
	"net/http"

	"github.com/apex/gateway"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	loadedConfig "github.com/adamjq/serverless-twirp/internal/config"
	"github.com/adamjq/serverless-twirp/internal/server"
	"github.com/adamjq/serverless-twirp/internal/stores"
	"github.com/adamjq/serverless-twirp/internal/userpb"
)

func main() {
	cfg := loadedConfig.NewConfig()

	// ddbOverride := "http://localstack:4566"

	// // this is necessary to override localstack endpoints in development mode
	// customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
	//   if service == dynamodb.ServiceID {
	// 		log.Printf("Overriding DynamoDB endpoint to %s...", ddbOverride)
	// 		return aws.Endpoint{
	// 			PartitionID:   "aws",
	// 			URL:           ddbOverride,
	// 			SigningRegion: "us-east-1",
	// 		}, nil
	//   }
	//   // returning EndpointNotFoundError will allow the service to fallback to it's default resolution
	//   return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	// })
	//awsCfg, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithEndpointResolver(customResolver))

	awsCfg, err := awsConfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := dynamodb.NewFromConfig(awsCfg)

	userStore := stores.NewUserStore(cfg.BackendTable, svc)
	s := server.NewServer(userStore)

	twirpHandler := userpb.NewUserServiceServer(s)

	mux := http.NewServeMux()
	mux.Handle(userpb.UserServicePathPrefix, twirpHandler)

	log.Fatal(gateway.ListenAndServe("", mux))
}
