package http

import (
	atDomain "github.com/JCFlores93/bookstore_oauth_api/src/domain/access_token"
	"github.com/JCFlores93/bookstore_oauth_api/src/services/access_token"
	"github.com/JCFlores93/bookstore_oauth_api/src/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type accessTokenHandler struct {
	service access_token.Service
}

type AccessTokenHandler interface {
	GetById(c *gin.Context)
	Create(c *gin.Context)
}

func NewAccessTokenHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (handler *accessTokenHandler) GetById(c *gin.Context) {
	accessToken, err := handler.service.GetById(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (handler *accessTokenHandler) Create(c *gin.Context) {
	var at atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	accessToken, err := handler.service.Create(at)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusCreated, accessToken)
}