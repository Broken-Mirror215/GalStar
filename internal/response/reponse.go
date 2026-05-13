package response

import "github.com/gin-gonic/gin"

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"` //这个？为什么要写接口类型？
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Body{
		Code:    0, //这个是业务成功了
		Message: "success",
		Data:    data,
	})
}

func Fail(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Body{
		Message: message,
		Data:    nil,
		Code:    code,
	})
}
