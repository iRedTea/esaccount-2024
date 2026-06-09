package handler

import (
	"esaccount/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "UserInfo"
	userHeader          = "UserHeader"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		util.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	user, err := h.services.Authorization.Authorize(header)
	if err != nil {
		util.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	_, err = h.repo.GetAccount(user.Id)
	if err != nil {
		h.repo.CreateAccount(user.Id)
	}

	c.Set(userCtx, user.Id)
	c.Set(userHeader, header)
}

