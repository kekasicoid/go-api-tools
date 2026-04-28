package model

import (
	"fmt"

	"github.com/gin-gonic/gin"
	response "github.com/kekasicoid/kekasigohelper/dto/response/v1"
)

func resolveRequestID(c *gin.Context) string {
	if v, ok := c.Get(HeadRequestIDKey); ok {
		if ctxID, ok := v.(string); ok && ctxID != "" {
			return ctxID
		}
	}

	return c.GetHeader(HeadRequestIDKey)
}

func RespBadRequest(c *gin.Context, desc string) {
	httpCode := 400
	resp := response.New(resolveRequestID(c))
	resp.Code = fmt.Sprintf("%d", httpCode)
	resp.Desc = desc
	c.JSON(httpCode, resp)
}

func RespSuccess(c *gin.Context, desc string, data interface{}) {
	httpCode := 200
	resp := response.New(resolveRequestID(c))
	resp.Code = fmt.Sprintf("%d", httpCode)
	resp.Desc = desc
	resp.Data = data
	c.JSON(httpCode, resp)
}
