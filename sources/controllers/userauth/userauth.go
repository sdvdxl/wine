package userauth

import (
	"fmt"
	gincache "github.com/gin-gonic/contrib/cache"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/bean"
	. "github.com/sdvdxl/wine/sources/util/cache"
	"github.com/sdvdxl/wine/sources/util/config"
	"github.com/sdvdxl/wine/sources/util/constant"
	. "github.com/sdvdxl/wine/sources/util/log"
	"github.com/sdvdxl/wine/sources/util/stringutils"
	"net/http"
	"strings"
)

type userAuth struct {
	context *gin.Context
}

func Auth(c *gin.Context) *userAuth {
	return &userAuth{c}
}

func (auth userAuth) IsLogined() bool {
	if auth.CurUser() == nil {
		return false
	}

	return true
}

func (auth *userAuth) IsRole(role string) bool {
	curUser := auth.CurUser()
	return strings.ToLower(curUser.Role) == strings.ToLower(role)
}

func (auth userAuth) isRequiredLogin() bool {
	path := auth.context.Request.URL.Path
	if path == "/user/login.html" || path == "/user/regist.html" {
		return false
	}

	for _, v := range config.AuthPageConfig.LoginPaths.Values() {
		Logger.Debug("check requied login pages, path :%v, check path:%v", path, v)
		if strings.HasPrefix(path, fmt.Sprint(v)) {
			return true
		}
	}

	if config.AuthPageConfig.LoginPages.Contains(path) {
		return true
	}
	return false
}

//退出
func (auth *userAuth) Logout() {
	session := sessions.Default(auth.context)
	Cache.Delete(stringutils.ToString(session.Get(constant.LOGIN_PID)))
	session.Delete(constant.LOGIN_PID)
	session.Save()
}

//获取当前登录用户，如果从cache中有错误，或者pid cooke已经失效，则返回nil，否则返回当前登录的用户信息
func (auth *userAuth) CurUser() *bean.User {
	session := sessions.Default(auth.context)
	pid := session.Get(constant.LOGIN_PID)

	pidStr := fmt.Sprintf("%v", pid)
	user := &bean.User{}
	if err := Cache.Get(pidStr, user); err != nil {
		if err == gincache.ErrCacheMiss {
			Logger.Debug("user info not found in cache pid: %v", pidStr)
		} else {
			Logger.Error("error occured when get login pid from cache: %v", err)
		}
		return nil
	}

	return user

}

func CheckLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := Auth(c)
		if auth.IsLogined() || !strings.HasSuffix(c.Request.URL.Path, ".html") {
			c.Next()
			return
		}

		Logger.Debug("checking login page,request page:%v", c.Request.URL.Path)

		//如果没有登录，且该页面是要求登录才能访问，那么跳转到登录页面
		if !auth.IsLogined() && auth.isRequiredLogin() {
			Logger.Debug("user not logined")
			Logger.Debug("required login")
			c.Redirect(http.StatusMovedPermanently, "/user/login.html")
		}
	}

}
