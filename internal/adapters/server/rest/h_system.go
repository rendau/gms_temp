package rest

import "github.com/gin-gonic/gin"

func (o *St) hSystemCronTick5m(c *gin.Context) {
	go o.ucs.SystemCronTick5m()
}

func (o *St) hSystemCronTick15m(c *gin.Context) {
	go o.ucs.SystemCronTick15m()
}

func (o *St) hSystemCronTick30m(c *gin.Context) {
	go o.ucs.SystemCronTick30m()
}
