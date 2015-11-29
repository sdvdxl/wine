package userinfo

import (
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/bean"
	"github.com/sdvdxl/wine/sources/controllers/userauth"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/wine/sources/util/db"
	. "github.com/sdvdxl/wine/sources/util/log"
	"net/http"
	"regexp"
)

//手机号是否可用
func IsPhoneValid(c *gin.Context) {
	phone := c.Query("phone")
	Logger.Info(phone)
	user := &bean.User{Phone: phone}
	userCount, err := db.Engine.Count(user)
	util.LogError(err)
	Logger.Info(userCount)

	if !IsLegalPhoneNumber(phone) {
		c.JSON(http.StatusOK, util.JsonResult{Success: true, Msg: "请输入正确的手机号(11位)"})
		return
	}

	if userCount > 0 {
		c.JSON(http.StatusOK, util.JsonResult{Success: true, Msg: "该手机号已经被注册"})
	} else {
		c.JSON(http.StatusOK, util.JsonResult{Success: true, Msg: "手机号可用"})
	}
}

func IsLegalPhoneNumber(phone string) bool {
	regex := regexp.MustCompile("^\\d{11}$")
	if regex.MatchString(phone) {
		return true
	}

	return false
}

//校验是否是合法用户名，4-64个，字母和数字
func IsLegalUsername(username string) bool {
	regex := regexp.MustCompile("^[\\w\\d]{4,64}$")
	if regex.MatchString(username) {
		return true
	}

	return false
}

func GetUserInfoHandler(c *gin.Context) {
	auth := userauth.Auth(c)
	if auth.IsLogined() {
		c.JSON(http.StatusOK, util.JsonResult{Success: true, Data: auth.CurUser()})
	}
}
