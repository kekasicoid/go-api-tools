package model

import (
	"fmt"

	"github.com/gin-gonic/gin"
	response "github.com/kekasicoid/kekasigohelper/dto/response/v1"
)

func RespBadRequest(c *gin.Context, id string, desc string) {
	httpCode := 400
	resp := response.New(id)
	resp.Code = fmt.Sprintf("%d", httpCode)
	resp.Desc = desc
	c.JSON(httpCode, resp)
}

func RespSuccess(c *gin.Context, id string, desc string, data interface{}) {
	httpCode := 200
	resp := response.New(id)
	resp.Code = fmt.Sprintf("%d", httpCode)
	resp.Desc = desc
	resp.Data = data
	c.JSON(httpCode, resp)
}
