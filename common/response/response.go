package response

import (
	"bifrost/common"
	"bifrost/common/errorx"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Response 封装成功与业务错误返回
func Response(c *gin.Context, resp interface{}, err error, code ...int) {
	if err != nil {
		RespErr(c, err)
		return
	}

	// 设置原始响应 Header
	if resp != nil {
		respByte, _ := json.Marshal(resp)
		c.Header(common.RespBody, string(respByte))
	}

	c.Header(common.StatusCode, strconv.Itoa(http.StatusOK))

	if len(code) > 0 && errorx.ErrMap[code[0]] != nil {
		c.JSON(http.StatusOK, Body{
			Code: code[0],
			Msg:  errorx.ErrMap[code[0]].Error(),
			Data: resp,
		})
	} else {
		c.JSON(http.StatusOK, Body{
			Code: errorx.Success,
			Msg:  "OK",
			Data: resp,
		})
	}
}

// RespErr 封装错误响应
func RespErr(c *gin.Context, err error) {
	code, err := errorx.TransCodeErr(err)

	if _, ok := err.(*errorx.CodeError); !ok {
		c.Header(common.StatusCode, strconv.Itoa(http.StatusInternalServerError))
		c.Header(common.RespErr, fmt.Sprintf("%+v", err))
		err = errorx.NewCodeError(errorx.DefaultCode)
		code = errorx.DefaultCode
	} else {
		c.Header(common.StatusCode, strconv.Itoa(http.StatusOK))
	}

	c.JSON(http.StatusOK, Body{
		Code: code,
		Msg:  err.Error(),
	})
}
