package stores

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	NotFoundError string = "NotFoundError"
)

type Users interface {
	GetUser(ctx context.Context, userId, orgId string) (*User, error)
	StoreUser(ctx context.Context, user StoreUser) (*string, error)
}

type UserStore struct {
	Users
	ddb       *dynamodb.Client
	tableName string
}

func NewUserStore(tableName string, ddb *dynamodb.Client) *UserStore {
	return &UserStore{
		tableName: tableName,
		ddb:       ddb,
	}
}

type User struct {
	PK             string    `validate:"required"`
	SK             string    `validate:"required"`
	UserID         string    `validate:"required"`
	OrganisationID string    `validate:"required"`
	FirstName      string    `validate:"required"`
	LastName       string    `validate:"required"`
	Role           string    `validate:"required"`
	CreatedTime    time.Time `validate:"required"`
}

type StoreUser struct {
	OrganisationID string `validate:"required"`
	FirstName      string `validate:"required"`
	LastName       string `validate:"required"`
	Role           string `validate:"required"`
}

// GetUser fetches and returns a user stored in DynamoDB
func (us *UserStore) GetUser(ctx context.Context, orgId, userId string) (*User, error) {
	pk := formatPK(orgId)
	sk := formatSK(userId)

	result, err := us.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(us.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})

	if err != nil {
		log.Printf("ERROR: %+v", err)
		return nil, errors.New("error fetching item from DynamoDB")
	}
	if result.Item == nil {
		log.Printf("ERROR: No user found in DynamoDB for OrgId: %s UserId: %s", orgId, userId)
		return nil, errors.New(fmt.Sprintf("%s: couldn't find item in DynamoDB", NotFoundError))
	}

	user := &User{}
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, errors.New("failed to unmarshall user")
	}

	return user, nil
}

// StoreUser stores a user in DynamoDB and returns their UserID
func (us *UserStore) StoreUser(ctx context.Context, user StoreUser) (*string, error) {
	newUserId := uuid.New().String()
	pk := formatPK(user.OrganisationID)
	sk := formatSK(newUserId)
	createdTime := time.Now().UTC()

	userToStore := User{
		PK:             pk,
		SK:             sk,
		UserID:         newUserId,
		OrganisationID: user.OrganisationID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Role:           user.Role,
		CreatedTime:    createdTime,
	}

	item, err := attributevalue.MarshalMap(userToStore)
	if err != nil {
		log.Printf("ERROR: %+v", err)
		return nil, errors.New("failed to marshall user")
	}

	_, err = us.ddb.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(us.tableName),
		Item:      item,
	})

	if err != nil {
		log.Printf("ERROR: %+v", err)
		return nil, errors.New("error fetching item from DynamoDB")
	}

	return &newUserId, nil
}

func formatPK(orgId string) string {
	return fmt.Sprintf("ORG#%s", orgId)
}

func formatSK(userId string) string {
	return fmt.Sprintf("USER#%s", userId)
}
