package route

import (
	"nats-go/server/config"
	"nats-go/server/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()

	// gin 檔案上傳body限制
	r.MaxMultipartMemory = 64 << 20

	// Middleware
	//r.Use(middleware.LogRequest())
	r.Use(middleware.ErrorResponse())

	corsConf := cors.DefaultConfig()
	corsConf.AllowCredentials = true
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	corsConf.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "Bearer", "Accept-Language"}
	corsConf.AllowOriginFunc = config.GetCorsRule
	r.Use(cors.New(corsConf))

	apiRoute := r.Group("/api")

	AgentV1(apiRoute)
	TaskV1(apiRoute)

	return r
}
