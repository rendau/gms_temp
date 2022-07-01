package rest

import (
	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/gms_temp/internal/domain/entities"
)

// @Router   /config [put]
// @Tags     config
// @Summary  Update configs
// @Accept   json
// @Param    body  body  entities.ConfigSt  false  "body"
// @Success  200
// @Failure  400  {object}  dopTypes.ErrRep
func (o *St) hConfigUpdate(c *gin.Context) {
	reqObj := &entities.ConfigSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	dopHttps.Error(c, o.ucs.ConfigSet(o.getRequestContext(c), reqObj))
}
