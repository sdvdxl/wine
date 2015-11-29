package index

import (
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/controllers/userauth"
	. "github.com/sdvdxl/wine/sources/util/log"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	auth := userauth.Auth(c)
	if auth.IsLogined() {
		Logger.Debug("current login user is: %v", auth.CurUser())
		c.Redirect(http.StatusMovedPermanently, "/order/order_list.html")
	} else {
		c.Redirect(http.StatusMovedPermanently, "/user/login.html")
	}
}
