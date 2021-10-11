package server

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/adamjq/serverless-twirp/internal/stores"
	"github.com/adamjq/serverless-twirp/internal/userpb"
	"github.com/google/uuid"
	"github.com/twitchtv/twirp"
)

// Server implements the Hello service
type Server struct {
	userStore stores.Users
}

func NewServer(userStore stores.Users) *Server {
	return &Server{
		userStore: userStore,
	}
}

func (s *Server) GetUser(ctx context.Context, in *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	orgId := in.GetOrganisationId()
	userId := in.GetUserId()

	if orgId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "organisation_id is required")
	}
	if userId == "" {
		return nil, twirp.NewError(twirp.InvalidArgument, "user_id is required")
	}

	user, err := s.userStore.GetUser(ctx, orgId, userId)

	if err != nil {
		if strings.Contains(err.Error(), stores.NotFoundError) {
			return nil, twirp.NotFoundError(fmt.Sprintf("User %s not found for Organisation %s", userId, orgId))
		}

		return nil, twirp.WrapError(twirp.NewError(twirp.Internal, "something went wrong"), err)
	}

	return &userpb.GetUserResponse{
		User: mapUserToTwirp(user),
	}, nil
}

func (s *Server) StoreUser(ctx context.Context, in *userpb.StoreUserRequest) (*userpb.StoreUserResponse, error) {

	newUserId := uuid.New().String()

	storeUser := stores.StoreUser{
		OrganisationID: in.GetOrganisationId(),
		FirstName:      in.GetFirstName(),
		LastName:       in.GetLastName(),
		Role:           in.GetRole().Enum().String(),
	}
	log.Printf("SERVER: Storing user %+v", storeUser)

	user, err := s.userStore.StoreUser(ctx, newUserId, storeUser)

	if err != nil {
		return nil, twirp.WrapError(twirp.NewError(twirp.Internal, "something went wrong"), err)
	}

	twirpUser := mapUserToTwirp(user)
	log.Printf("SERVER: Returning stored user %+v", twirpUser)

	response := userpb.StoreUserResponse{
		User: twirpUser,
	}

	return &response, nil
}

func mapUserToTwirp(user *stores.User) *userpb.User {
	return &userpb.User{
		OrganisationId: user.OrganisationID,
		UserId:         user.UserID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Role:           mapUserRole(user.Role),
	}
}

func mapUserRole(role string) userpb.UserRole {
	if role == "USER_ROLE_READONLY" {
		return 1
	}
	if role == "USER_ROLE_ADMIN" {
		return 2
	}
	return 0
}
