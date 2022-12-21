package v1

import (
	"errors"
	"net/http"

	"github.com/MuhammadyusufAdhamov/blog_project/pkg/utils"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationPayLoadKey = "authorization_payload"
)

func (h *handlerV1) AuthMiddleWare(c *gin.Context) {
	accesToken := c.GetHeader(authorizationHeaderKey)

	if len(accesToken) == 0 {
		err := errors.New("authorization header is not provided")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	payload, err := utils.VerifyToken(accesToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.Set(authorizationPayLoadKey, payload)
	c.Next()
}

func (m *handlerV1) GetAuthPayload(ctx *gin.Context) (*utils.Payload, error) {
	i, exists := ctx.Get(authorizationPayLoadKey)
	if !exists {
		return nil, errors.New("")
	}

	payload, ok := i.(*utils.Payload)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return payload, nil
}