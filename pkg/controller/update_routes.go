package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/litepubl/test-treasury/pkg/importer"
)

type (
	UpdateRoutes struct {
		updater importer.DataUpdater
	}

	updateResponse struct {
		Result bool   `json:"result"`
		Info   string `json:"info"`
		Code   int    `json:"code"`
	}

	stateResponse struct {
		Result bool   `json:"result"`
		Info   string `json:"info"`
	}
)

func NewUpdateRoutes(u importer.DataUpdater) *UpdateRoutes {
	return &UpdateRoutes{
		updater: u,
	}
}

func (u *UpdateRoutes) Update(c *gin.Context) {
	err := u.updater.Update(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, updateResponse{
			Result: false,
			Info:   "service unavailable",
			Code:   503,
		})

		return
	}

	c.JSON(http.StatusOK, updateResponse{
		Result: true,
		Info:   "",
		Code:   200,
	})
}

func (u *UpdateRoutes) State(c *gin.Context) {
	state := u.updater.State()
	c.JSON(http.StatusOK, stateResponse{
		Result: state == importer.Ok,
		Info:   state.String(),
	})
}
