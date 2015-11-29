package login

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/bean"
	"github.com/sdvdxl/wine/sources/controllers/userauth"
	"github.com/sdvdxl/wine/sources/util"
	. "github.com/sdvdxl/wine/sources/util/cache"
	"github.com/sdvdxl/wine/sources/util/constant"
	"github.com/sdvdxl/wine/sources/util/db"
	. "github.com/sdvdxl/wine/sources/util/log"
	"github.com/sdvdxl/wine/sources/util/stringutils"
	"net/http"
	"strings"
	"time"
)

func init() {
	captcha.NewMemoryStore(1000, time.Minute*5)
}

// 处理登陆
func LoginHandler(c *gin.Context) {
	session := sessions.Default(c)

	//校验验证码

	loginCaptcha := c.PostForm("captcha")
	Logger.Debug("login captcha is:%v", loginCaptcha)
	memCaptcha := stringutils.ToString(session.Get(constant.LOGIN_CAPTCHA))
	Logger.Debug("mem captcha:%v", memCaptcha)
	if !strings.EqualFold(memCaptcha, loginCaptcha) {
		session.Set(constant.LOGIN_ERROR_MSG, "验证码不正确")
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/user/login.html")
		return
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) == 0 || len(password) == 0 {
		session.Set(constant.LOGIN_ERROR_MSG, "请填写用户名和密码")
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/user/login.html")
		return
	}

	user := &bean.User{Phone: username}
	found, err := db.Engine.Get(user)
	util.PanicError(err)

	//校验用户名和密码
	if !found || (util.HashAndSalt(password, user.Salt) != user.Password) {
		Logger.Error("user %v  not found in user", username)

		session.Set(constant.LOGIN_ERROR_MSG, "用户名或者密码错误")
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/user/login.html")
		return
	}

	SetLoginState(*user, c)

	c.Redirect(http.StatusMovedPermanently, "/order/order_list.html")
}

func TemplateHandler(c *gin.Context) {
	session := sessions.Default(c)
	userAuth := userauth.Auth(c)
	Logger.Debug("user has logined: %v", userAuth.IsLogined())

	if userAuth.IsLogined() {
		c.Redirect(http.StatusMovedPermanently, "/")
	}

	loginErrorMsg := session.Get(constant.LOGIN_ERROR_MSG)
	session.Delete(constant.LOGIN_ERROR_MSG)
	session.Save()

	c.HTML(http.StatusOK, "login.tmpl", gin.H{"loginErrorMsg": loginErrorMsg})

}

func LogoutHandler(c *gin.Context) {
	auth := userauth.Auth(c)
	auth.Logout()
	c.Redirect(http.StatusMovedPermanently, "/")
}

//需要user的phone，email
func SetLoginState(user bean.User, c *gin.Context) {
	session := sessions.Default(c)
	timestamp := time.Now().UnixNano()
	pid := util.HashPid(user.Phone, user.Email, timestamp)
	session.Set(constant.LOGIN_PID, pid)
	session.Set(constant.LOGIN_TIMESTAMP, timestamp)
	if err := Cache.Set(pid, user, constant.LOGIN_EXPIRED_TIME); err != nil {
		Logger.Error("error occured when set user to cache, %v", err)
	}

	Logger.Debug("login pid : %v", pid)
	if err := session.Save(); err != nil {
		Logger.Error("error occured when save session, %v", err)
	}
}
