package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	target *url.URL
}

func NewController() *Controller {
	return &Controller{}
}
func (c *Controller) Proxy() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(c.target)
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
