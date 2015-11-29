// main.go
package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/controllers/captcha"
	"github.com/sdvdxl/wine/sources/controllers/index"
	"github.com/sdvdxl/wine/sources/controllers/login"
	"github.com/sdvdxl/wine/sources/controllers/nav"
	"github.com/sdvdxl/wine/sources/controllers/order"
	"github.com/sdvdxl/wine/sources/controllers/regist"
	"github.com/sdvdxl/wine/sources/controllers/user"
	"github.com/sdvdxl/wine/sources/controllers/userauth"
	"github.com/sdvdxl/wine/sources/controllers/userinfo"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/wine/sources/util/db"
	"github.com/sdvdxl/wine/sources/util/log"
	"github.com/sdvdxl/wine/sources/util/mail"
	"github.com/sdvdxl/wine/sources/util/render"
	"github.com/sdvdxl/wine/sources/util/session"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	g *gin.Engine
)

func init() {
	//初始化gin处理引擎
	//	gin.SetMode(gin.ReleaseMode)
	g = gin.New()
	g.Use(HandlerError())

	{
		funcMap := template.FuncMap{"Equals": func(v1, v2 interface{}) bool {
			log.Logger.Debug("invoke function Equals")
			return v1 == v2
		}}
		tmp := template.New("myTemplate")
		templatePages := TemplatesFinder("templates")
		tmp.Funcs(funcMap).ParseFiles(templatePages...)

		g.SetHTMLTemplate(tmp)

		{ //这三个顺序不能变更,否则得不到正常处理
			//先设置/读取session信息
			g.Use(sessions.Sessions("my_session", session.SessionStore))

			//然后校验请求的URL
			g.Use(userauth.CheckLoginPage())

			//最后处理静态文件
			g.Use(static.ServeRoot("/", "static")) // static files have higher priority over dynamic routes

		}

	}
}

func HandlerError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 10240)
				runtime.Stack(trace, true)
				log.Logger.Error("%s, \n%s", err, trace)

				if strings.HasSuffix(c.Request.URL.Path, ".html") {
					c.HTML(http.StatusInternalServerError, "500.tmpl", nil)
				} else {
					r := render.New(c)
					r.JSON(util.JsonResult{Msg: "系统错误"})
				}
			}
		}()
		c.Next()
	}
}

func main() {
	defer log.Logger.Close()
	defer db.Close()

	log.Logger.Info("starting server...")

	//==================================   nav top sidebar  ======================================
	{
		g.Any("/sidebar.html", nav.SidebarTemplateHandler)
		g.Any("/header.html", nav.HeaderTemplateHandler)
		g.Any("/footer.html", nav.FooterTemplateHandler)
	}

	//==================================   404, index  ======================================
	{
		//404
		g.NoRoute(func(c *gin.Context) {
			log.Logger.Debug("page [%v] not found, redirect to 404.html", c.Request.URL.Path)
			if strings.HasSuffix(c.Request.URL.Path, ".html") {
				c.HTML(http.StatusNotFound, "404.tmpl", nil)
			} else {
				r := render.New(c)
				r.JSON(util.JsonResult{Msg: "请求资源不存在"})
			}
		})

		//首页
		g.Any("/", index.IndexHandler)
	}

	//==================================  用户相关  ======================================
	{
		userGroup := g.Group("/user")
		{
			{ //处理登陆
				userGroup.Any("/logout", login.LogoutHandler)
				userGroup.GET("/getUserInfo", userinfo.GetUserInfoHandler)
				userGroup.POST("/login", login.LoginHandler)
				userGroup.GET("/login.html", login.TemplateHandler)
				userGroup.GET("/getImageCaptcha", captcha.GetImageCaptcha)

				userGroup.GET("/user_list", user.UserListHandler)
				userGroup.POST("/user_add", user.UserAddHandler)
				userGroup.POST("/user_delete", user.UserDeleteHandler)
				userGroup.POST("/user_update", user.UserUpdateHandler)
				userGroup.POST("/change_password", user.ChangePasswordHandler)
			}

			{
				userGroup.GET("/me", user.MyProfileHandler)
			}

			{ //注册
				userGroup.POST("/regist", regist.RegistHandler)
				userGroup.POST("/registGuid", regist.RegistHandler)
				userGroup.GET("/isPhoneValid", userinfo.IsPhoneValid)
				userGroup.GET("/getPhoneCaptcha", captcha.GetPhoneCaptcha)
			}
		}
	}

	//==================================  订单  ======================================
	{
		orderGroup := g.Group("/order")
		orderGroup.GET("/order_list", order.OrderListHandler)
		orderGroup.POST("/order_update", order.OrderUpdateHandler)
		orderGroup.POST("/order_delete", order.OrderDeleteHandler)
		orderGroup.POST("/order_add_comments", order.OrderAddCommentsHandler)
		orderGroup.POST("/order_upload", order.OrderUploadHandler)
		orderGroup.GET("/order_upload.html", order.OrderUploadTemplateHandler)
		orderGroup.GET("/order_comments", order.OrderCommentsHandler) //添加订单注释
	}

	//==================================  mail  ======================================
	{
		g.GET("/mail", func(c *gin.Context) {

			mail.SendEmail([]string{c.Query("to")}, "test", "mail.tmpl", c.Query("content"))

			c.JSON(http.StatusOK, util.JsonResult{Msg: "发送邮件成功"})

		})
	}

	log.Logger.Info("server started ")
	g.Run(":8085")
}

func TemplatesFinder(templateDirName string) []string {
	templatePages := make([]string, 0, 10)
	filepath.Walk(templateDirName, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			templatePages = append(templatePages, path)
		}
		return nil
	})

	log.Logger.Debug("templates:%v", templatePages)
	return templatePages
}
