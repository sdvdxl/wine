package regist

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/sdvdxl/go-uuid/uuid"
	mdl "github.com/sdvdxl/wine/sources/bean"
	"github.com/sdvdxl/wine/sources/controllers/login"
	"github.com/sdvdxl/wine/sources/controllers/userinfo"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/wine/sources/util/constant"
	"github.com/sdvdxl/wine/sources/util/db"
	. "github.com/sdvdxl/wine/sources/util/log"
	"github.com/sdvdxl/wine/sources/util/stringutils"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//普通用户（游客）注册
func RegistHandler(c *gin.Context) {
	session := sessions.Default(c)
	phone := c.PostForm("phone")
	Logger.Info("register phone :%v", phone)
	phone = strings.TrimSpace(phone)
	if !userinfo.IsLegalPhoneNumber(phone) {
		render.WriteJSON(c.Writer, "请填写正确的手机号(11位)")
		return
	}

	if !strings.EqualFold(phone, stringutils.ToString(session.Get(constant.PHONE_NUMBER))) {
		c.JSON(http.StatusOK, util.JsonResult{Success: false, Msg: "下发验证码的手机号不是当前手机号"})
		return
	}

	if !strings.EqualFold(strings.ToLower(stringutils.ToString(session.Get(constant.PHONE_CAPTCHA))), strings.ToLower(c.PostForm("captcha"))) {
		c.JSON(http.StatusOK, util.JsonResult{Success: false, Msg: "验证码不正确"})
		return

		lastPhoneCaptcha := session.Get(constant.PHONE_CAPTCHA_LAST)
		lastPhoneCaptchaTimestamp, _ := lastPhoneCaptcha.(int64)
		timeDifference := (time.Now().UnixNano() - lastPhoneCaptchaTimestamp) / (1000 * 1000 * 1000 * 60)
		if timeDifference > constant.PHONE_CAPTCHA_EXPIRED_MINUTES {
			c.JSON(http.StatusOK, util.JsonResult{Success: false, Msg: fmt.Sprintf("验证码已失效，请在验证码下发后%v分钟内提交", constant.PHONE_CAPTCHA_EXPIRED_MINUTES)})
			return
		}
	}

	user := &mdl.User{Phone: phone}
	userCount, err := db.Engine.Count(user)
	util.LogError(err)

	if userCount != 0 {
		render.WriteJSON(c.Writer, util.JsonResult{Success: true, Msg: "该手机号已经被注册"})
		return
	}

	user.Salt = fmt.Sprintf("%v", rand.New(rand.NewSource(time.Now().UnixNano())).Float64())
	user.Password = c.PostForm("password")
	if len(user.Password) < 6 {
		c.JSON(http.StatusOK, util.JsonResult{Msg: "密码长度大于6"})
		return
	}
	user.Password = util.HashAndSalt(user.Password, user.Salt)
	user.Uuid = uuid.New()
	Logger.Info(user.Uuid)
	nickname := strings.TrimSpace(c.PostForm("nickname"))
	Logger.Debug("nickname: %v, size: %v", nickname, len(nickname))
	if len(nickname) == 0 {
		user.Nickname = user.Phone
	} else {
		user.Nickname = nickname
	}
	count, err := db.Engine.Insert(user)
	if err != nil || count == 0 {
		Logger.Error(err)
		c.JSON(http.StatusOK, util.JsonResult{Msg: "系统错误"})
		return
	}

	//注册成功，设置为登录
	login.SetLoginState(*user, c)

	c.JSON(http.StatusOK, util.JsonResult{Msg: "注册成功", Success: true})

}
