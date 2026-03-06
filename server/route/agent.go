package route

import (
	"nats-go/server/api"

	"github.com/gin-gonic/gin"
)

func AgentV1(r *gin.RouterGroup) {
	v1Route := r.Group("/v1/agent")

	v1Route.GET("", func(c *gin.Context) {
		api.ListAgent(c)
	})

}
