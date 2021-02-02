package proxy

import (
	"github.com/gin-gonic/gin"

	"github.com/searKing/sole/web/golang/app/configs/values"
)

func SetRouter(router gin.IRouter) gin.IRouter {
	proxy := NewController()
	router.Any(values.Proxy, proxy.Proxy())
	return router
}
