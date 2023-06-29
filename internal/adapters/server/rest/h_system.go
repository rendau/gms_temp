package rest

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (o *St) hSystemCronTick5m(c *gin.Context) {
	o.ucs.SystemCronTick5m(o.hSystemCronTickParseQueryT(c))
}

func (o *St) hSystemCronTick15m(c *gin.Context) {
	o.ucs.SystemCronTick15m(o.hSystemCronTickParseQueryT(c))
}

func (o *St) hSystemCronTick30m(c *gin.Context) {
	o.ucs.SystemCronTick30m(o.hSystemCronTickParseQueryT(c))
}

func (o *St) hSystemCronTickParseQueryT(c *gin.Context) time.Time {
	var t time.Time
	if tStr := c.Query("t"); tStr != "" {
		ts, err := time.Parse(time.RFC3339, tStr)
		if err == nil {
			t = ts
		}
	}
	return t
}
