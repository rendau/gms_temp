package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
)

// @Router      /dic [get]
// @Tags        dic
// @Summary     dictionaries
// @Description Get all dictionaries
// @Produce     json
// @Success     200 {object} entities.DicSt
// @Failure     400 {object} dopTypes.ErrRep
func (o *St) hDicGet(c *gin.Context) {
	result, err := o.ucs.DicGet(o.getRequestContext(c))
	if dopHttps.Error(c, err) {
		return
	}

	c.JSON(http.StatusOK, result)
}
