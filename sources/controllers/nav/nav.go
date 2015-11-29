package nav

import (
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/controllers/userauth"
	"github.com/sdvdxl/wine/sources/util/log"
	"net/http"
	"strings"
)

func SidebarTemplateHandler(c *gin.Context) {
	var curNav string
	referPath := c.Request.Referer()
	referPath = referPath[strings.Index(referPath, c.Request.Host)+len(c.Request.Host):]

	log.Logger.Debug("host:%v, sidebar refer is: %v", c.Request.Host, referPath)
	if strings.HasPrefix(referPath, "/order") {
		curNav = "order"
	} else if strings.HasPrefix(referPath, "/user") {
		curNav = "user"
	}
	auth := userauth.Auth(c)
	c.HTML(http.StatusOK, "sidebar.tmpl", gin.H{"curNav": curNav, "curUser": auth.CurUser()})
}

func HeaderTemplateHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "header.tmpl", nil)
}

func FooterTemplateHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "footer.tmpl", nil)
}
