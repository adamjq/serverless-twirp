package server

import (
	"context"
	"errors"
	"testing"

	"github.com/adamjq/serverless-twirp/internal/stores"
	"github.com/adamjq/serverless-twirp/internal/userpb"
	"github.com/adamjq/serverless-twirp/mocks/storemock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestServer_GetUser(t *testing.T) {
	assert := require.New(t)

	type test struct {
		name                   string
		input                  *userpb.GetUserRequest
		mockStoreResponseData  *stores.User
		mockStoreResponseError error
		want                   *userpb.GetUserResponse
		wantErr                bool
	}
	tests := []test{
		{
			name: "GetUser request succeeds",
			input: &userpb.GetUserRequest{
				UserId:         "1234",
				OrganisationId: "5678",
			},
			mockStoreResponseData: &stores.User{
				UserID:         "1234",
				OrganisationID: "5678",
				FirstName:      "Adam",
				LastName:       "Quigley",
				Role:           "USER_ROLE_ADMIN",
			},
			mockStoreResponseError: nil,
			want: &userpb.GetUserResponse{
				User: &userpb.User{
					UserId:         "1234",
					OrganisationId: "5678",
					FirstName:      "Adam",
					LastName:       "Quigley",
					Role:           userpb.UserRole_USER_ROLE_ADMIN,
				},
			},
			wantErr: false,
		},
		{
			name: "GetUser user not found in store",
			input: &userpb.GetUserRequest{
				UserId:         "1234",
				OrganisationId: "5678",
			},
			mockStoreResponseData:  nil,
			mockStoreResponseError: errors.New("NotFoundError"),
			want:                   nil,
			wantErr:                true,
		},
		{
			name: "GetUser invaid input - no organisation ID",
			input: &userpb.GetUserRequest{
				UserId:         "1234",
			},
			mockStoreResponseData:  nil,
			mockStoreResponseError: nil,
			want:                   nil,
			wantErr:                true,
		},
		{
			name: "GetUser invaid input - no user ID",
			input: &userpb.GetUserRequest{
				OrganisationId: "5678",
			},
			mockStoreResponseData:  nil,
			mockStoreResponseError: nil,
			want:                   nil,
			wantErr:                true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			usersStoreMock := storemock.NewMockUsers(mockCtrl)

			s := NewServer(usersStoreMock)

			// only mock store calls if expected
			if tt.mockStoreResponseData != nil || tt.mockStoreResponseError != nil {
				usersStoreMock.EXPECT().GetUser(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(tt.mockStoreResponseData, tt.mockStoreResponseError)
			}

			ctx := context.Background()

			resp, err := s.GetUser(ctx, tt.input)

			if !tt.wantErr {
				assert.NoError(err)
				assert.Equal(resp, tt.want)
			} else {
				assert.Error(err)
			}
		})
	}
}

func TestServer_StoreUser(t *testing.T) {
	assert := require.New(t)

	type test struct {
		name                   string
		input                  *userpb.StoreUserRequest
		mockStoreResponseData  *stores.User
		mockStoreResponseError error
		want                   *userpb.StoreUserResponse
		wantErr                bool
	}
	tests := []test{
		{
			name: "StoreUser request succeeds",
			input: &userpb.StoreUserRequest{
				OrganisationId: "5678",
				FirstName: "Adam",
				LastName: "Quigley",
				Role: userpb.UserRole_USER_ROLE_ADMIN,
			},
			mockStoreResponseData: &stores.User{
				UserID:         "1234",
				OrganisationID: "5678",
				FirstName:      "Adam",
				LastName:       "Quigley",
				Role:           "USER_ROLE_ADMIN",
			},
			mockStoreResponseError: nil,
			want: &userpb.StoreUserResponse{
				User: &userpb.User{
					UserId:         "1234",
					OrganisationId: "5678",
					FirstName:      "Adam",
					LastName:       "Quigley",
					Role:           userpb.UserRole_USER_ROLE_ADMIN,
				},
			},
			wantErr: false,
		},
		{
			name: "StoreUser request fails",
			input: &userpb.StoreUserRequest{
				OrganisationId: "5678",
				FirstName: "Adam",
				LastName: "Quigley",
				Role: userpb.UserRole_USER_ROLE_ADMIN,
			},
			mockStoreResponseData: nil,
			mockStoreResponseError: errors.New("Internal error occurred"),
			want: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			usersStoreMock := storemock.NewMockUsers(mockCtrl)

			s := NewServer(usersStoreMock)

			// only mock store calls if expected
			if tt.mockStoreResponseData != nil || tt.mockStoreResponseError != nil {
				usersStoreMock.EXPECT().StoreUser(
					gomock.Any(), gomock.Any(),
				).Return(tt.mockStoreResponseData, tt.mockStoreResponseError)
			}

			ctx := context.Background()

			resp, err := s.StoreUser(ctx, tt.input)

			if !tt.wantErr {
				assert.NoError(err)
				assert.Equal(resp, tt.want)
			} else {
				assert.Error(err)
			}
		})
	}
}
