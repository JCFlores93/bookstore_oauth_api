package access_token

import (
	"github.com/JCFlores93/bookstore_oauth_api/src/domain/access_token"
	"github.com/JCFlores93/bookstore_oauth_api/src/repository/db"
	"github.com/JCFlores93/bookstore_oauth_api/src/repository/rest"
	"github.com/JCFlores93/bookstore_oauth_api/src/utils/errors"
	"strings"
)

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

type Service interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type Repository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}


func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) *service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, *errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, *errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// TODO: Support both grant types: client_credentials and password
	// Authenticate the user against the Users API:
	user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// Generate new AccessToken
	at := access_token.GetNewAccessToken(user.Id)
	at.Generate()

	// Save the new accesstoken in Cassandra
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) *errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}