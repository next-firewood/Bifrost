package errorx

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

const (
	Success     = 0
	DefaultCode = iota + 1000
	TokenExpired
	TokenRefresh
	NonPermission
)

var ErrMap = map[int]error{
	TokenExpired:  fmt.Errorf("令牌过期"),
	TokenRefresh:  fmt.Errorf("令牌刷新"),
	NonPermission: fmt.Errorf("无权限"),
}

// 微信相关错误
const (
	WxSuccess   = 0
	WxCommonErr = iota + 2000
)

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return e.Msg
}

func Is(err, target error) bool {
	e, ok := err.(*CodeError)
	if !ok {
		return false
	}

	et, ok := target.(*CodeError)
	if !ok {
		return false
	}

	if et.Code == e.Code && et.Msg == e.Msg {
		return true
	}

	return false
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

func NewCodeError(code int) error {
	return &CodeError{Code: code, Msg: ErrMap[code].Error()}
}

func NewRpcError(code int) error {
	return status.Error(codes.Code(code), ErrMap[code].Error())
}

func TransCodeErr(err error) (code int, res error) {
	stat, ok := status.FromError(err)
	if !ok {
		return code, err
	}

	code, _ = strconv.Atoi(stat.Code().String()[5 : len(stat.Code().String())-1])

	codeErr := ErrMap[code]

	if codeErr != nil && codeErr.Error() == stat.Message() {
		return code, NewCodeError(code)
	}

	return code, err
}
