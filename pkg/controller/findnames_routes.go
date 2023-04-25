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

	strong := "weak" != strings.ToLower(c.Query("type"))
	entities, err := f.finder.GetNames(c.Request.Context(), name, strong)
	if err != nil {
		//r.l.Error(err, "http - v1 - history")
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, "service unavailable")

		return
	}

	c.JSON(http.StatusOK, entities)
}
