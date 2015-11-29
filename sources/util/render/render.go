package render

import (
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/go-tools/json"
	"net/http"
)

type rend struct {
	context *gin.Context
}

type JSONRender struct {
}

func New(c *gin.Context) *rend {
	return &rend{c}
}

func (r *rend) JSON(data interface{}) {
	result, err := jsonutils.Marshal(data, true)
	util.PanicError(err)
	r.context.String(http.StatusOK, "%s", result)
}
