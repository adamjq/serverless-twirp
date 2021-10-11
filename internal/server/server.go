package server

import (
	"context"
	"log"

	"github.com/adamjq/serverless-twirp/internal/stores"
	"github.com/adamjq/serverless-twirp/internal/userpb"
	"github.com/google/uuid"
	"github.com/twitchtv/twirp"
)

// Server implements the Hello service
type Server struct {
	userStore *stores.UserStore
}

func NewServer(userStore *stores.UserStore) *Server {
	return &Server{
		userStore: userStore,
	}
}

func (s *Server) GetUser(ctx context.Context, in *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	orgId := in.GetOrganisationId()
	userId := in.GetUserId()

	log.Printf("SERVER: Fetching user OrgID: %v UserID %v from store", orgId, userId)

	user, err := s.userStore.GetUser(ctx, orgId, userId)

	if err != nil {

		// TODO: check not found case

		return nil, twirp.WrapError(twirp.NewError(twirp.Internal, "something went wrong"), err)
	}

	log.Printf("SERVER: Retrieved user from store %+v", user)

	twirpUser := mapUserToTwirp(user)
	log.Printf("SERVER: Returning fetched user %+v", twirpUser)

	return &userpb.GetUserResponse{
		User: twirpUser,
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
