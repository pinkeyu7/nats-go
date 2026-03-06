package api

import (
	"nats-go/pkg/er"
	"nats-go/server/dto/req"
	"nats-go/server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendAsyncTask(c *gin.Context) {
	t := &req.Task{}
	if err := c.ShouldBindJSON(t); err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// assert
	apiEnv := GetEnv()
	ts, err := service.NewTaskService(apiEnv.GetNC(), apiEnv.GetJS())
	if err != nil {
		serviceErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, err.Error(), err)
		_ = c.Error(serviceErr)
		return
	}

	// act
	res, err := ts.SendAsyncTask(t)
	if err != nil {
		serviceErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, err.Error(), err)
		_ = c.Error(serviceErr)
		return
	}

	c.JSON(http.StatusOK, res)
}

func SendSyncTask(c *gin.Context) {
	t := &req.Task{}
	if err := c.ShouldBindJSON(t); err != nil {
		paramErr := er.NewAppErr(http.StatusBadRequest, er.ErrorParamInvalid, err.Error(), err)
		_ = c.Error(paramErr)
		return
	}

	// assert
	apiEnv := GetEnv()
	ts, err := service.NewTaskService(apiEnv.GetNC(), apiEnv.GetJS())
	if err != nil {
		serviceErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, err.Error(), err)
		_ = c.Error(serviceErr)
		return
	}

	// act
	res, err := ts.SendSyncTask(t)
	if err != nil {
		serviceErr := er.NewAppErr(http.StatusInternalServerError, er.UnknownError, err.Error(), err)
		_ = c.Error(serviceErr)
		return
	}

	c.JSON(http.StatusOK, res)
}
