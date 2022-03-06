package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/apex/gateway"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	loadedConfig "github.com/adamjq/serverless-twirp/internal/config"
	"github.com/adamjq/serverless-twirp/internal/server"
	"github.com/adamjq/serverless-twirp/internal/stores"
	"github.com/adamjq/serverless-twirp/pkg/userpb"
)

func main() {
	cfg := loadedConfig.NewConfig()

	ddbClient, err := getDynamoDbClient(cfg)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	userStore := stores.NewUserStore(cfg.BackendTable, ddbClient)
	s := server.NewServer(userStore)

	twirpHandler := userpb.NewUserServiceServer(s)

	mux := http.NewServeMux()
	mux.Handle(userpb.UserServicePathPrefix, twirpHandler)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		bindAddr := fmt.Sprintf(":%s", cfg.Addr)
		log.Printf("Listening on %s", bindAddr)
		rpcServer := http.Server{
			Addr:         bindAddr,
			Handler:      mux,
			IdleTimeout:  20 * time.Second,
			ReadTimeout:  20 * time.Second,
			WriteTimeout: 20 * time.Second,
		}
		log.Fatal(rpcServer.ListenAndServe())
	} else {
		// running in lambda
		log.Fatal(gateway.ListenAndServe("", mux))
	}
}

func getDynamoDbClient(cfg *loadedConfig.Config) (*dynamodb.Client, error) {
	var ddbCfg aws.Config
	var err error
	if cfg.DynamoEndpoint != "" {
		ddbCfg, err = awsConfig.LoadDefaultConfig(context.Background(),
			awsConfig.WithRegion(cfg.AwsRegion),
			awsConfig.WithEndpointResolver(aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{URL: cfg.DynamoEndpoint}, nil
				})),
		)
	} else {
		ddbCfg, err = awsConfig.LoadDefaultConfig(context.Background())
	}
	if err != nil {
		return nil, err
	}

	ddbClient := dynamodb.NewFromConfig(ddbCfg)
	return ddbClient, err
}
