package route

import (
	"nats-go/server/api"

	"github.com/gin-gonic/gin"
)

func TaskV1(r *gin.RouterGroup) {
	v1Route := r.Group("/v1/task")

	v1Route.POST("async", func(c *gin.Context) {
		api.SendAsyncTask(c)
	})

	v1Route.POST("sync", func(c *gin.Context) {
		api.SendSyncTask(c)
	})
}
