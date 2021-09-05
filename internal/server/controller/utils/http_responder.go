package utils

import (
	"github.com/Viverov/guideliner/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpResponder interface {
	// Response create common webserver error.
	// title is general human-readable description of error, like "Validation error".
	// detail contains additional human-readable description, like "Some of required fields undefined".
	// devDetail add additional info only for dev or debug env.
	Response(ctx *gin.Context, statusCode int, title string, detail string, devDetail string)
	// InternalError is Response shortcut for 500 status code error
	InternalError(ctx *gin.Context, devDetail string)
}

type httpResponder struct {
	env string
}

func NewHttpResponder(env string) *httpResponder {
	return &httpResponder{env: env}
}

type ErrorResponse struct {
	ErrorTitle  string `json:"error_title,omitempty"`
	ErrorDetail string `json:"error_detail,omitempty"`
	DevDetail   string `json:"dev_detail"`
}

func (er *httpResponder) Response(ctx *gin.Context, statusCode int, title string, detail string, devDetail string) {
	res := ErrorResponse{
		ErrorTitle:  title,
		ErrorDetail: detail,
	}
	if er.env == config.EnvDebug || er.env == config.EnvDevelopment || er.env == config.EnvTest {
		res.DevDetail = devDetail
	}

	ctx.JSON(statusCode, res)
}

func (er *httpResponder) InternalError(ctx *gin.Context, devDetail string) {
	er.Response(ctx, http.StatusInternalServerError, "Something gone wrong", "", devDetail)
}
