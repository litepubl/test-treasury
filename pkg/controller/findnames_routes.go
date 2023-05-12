package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/litepubl/test-treasury/pkg/finder"
)

type FindnameRoutes struct {
	finder finder.DataFinder
}

func NewFindnameRoutes(f finder.DataFinder) *FindnameRoutes {
	return &FindnameRoutes{
		finder: f,
	}
}

func (f *FindnameRoutes) GetNames(c *gin.Context) {
	name := strings.TrimSpace(c.Query("name"))
	if name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Name argument not specified")
		return
	}

	strong := strings.ToLower(c.Query("type")) != "weak"
	entities, err := f.finder.Names(c.Request.Context(), name, strong)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, "service unavailable")

		return
	}

	c.JSON(http.StatusOK, entities)
}
