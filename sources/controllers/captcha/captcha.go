package captcha

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/wine/sources/bean"
	"github.com/sdvdxl/wine/sources/controllers/userinfo"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/wine/sources/util/constant"
	"github.com/sdvdxl/wine/sources/util/db"
	. "github.com/sdvdxl/wine/sources/util/log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func GetImageCaptcha(c *gin.Context) {
	session := sessions.Default(c)

	captchaId := captcha.New()
	bytes := captcha.RandomDigits(4)
	captcheString := util.GetCaptchaString(bytes)
	session.Set(constant.LOGIN_CAPTCHA, captcheString)
	Logger.Debug("new captcha:%v", captcheString)
	session.Save()
	image := captcha.NewImage(captchaId, bytes, captcha.StdWidth/2, captcha.StdHeight/2)
	_, err := image.WriteTo(c.Writer)
	if err != nil {
		Logger.Error("error occured when writing captcha image to response")
	}
}

//获取手机验证码
func GetPhoneCaptcha(c *gin.Context) {
	session := sessions.Default(c)
	phone := c.Query("phone")
	Logger.Info("get captcha phone:%v", phone)
	user := &bean.User{Phone: phone}
	userCount, err := db.Engine.Count(user)
	util.LogError(err)
	Logger.Info(userCount)

	if !userinfo.IsLegalPhoneNumber(phone) {
		c.JSON(http.StatusOK, util.JsonResult{Success: false, Msg: "请输入正确的手机号(11位)"})
		return
	}

	if userCount > 0 {
		c.JSON(http.StatusOK, util.JsonResult{Success: true, Msg: "该手机号已经被注册"})
		return
	}

	lastPhoneCaptcha := session.Get(constant.PHONE_CAPTCHA_LAST)
	lastPhoneCaptchaTimestamp, _ := lastPhoneCaptcha.(int64)
	timeDifference := (time.Now().UnixNano() - lastPhoneCaptchaTimestamp) / (1000 * 1000 * 1000)
	Logger.Debug("phone captcha time difference: %v", timeDifference)
	if timeDifference < constant.PHONE_CAPTCHA_SECONDS {
		c.JSON(http.StatusOK, util.JsonResult{Msg: fmt.Sprintf("还有%v秒可以再次发送验证码", int64(constant.PHONE_CAPTCHA_SECONDS)-timeDifference), Success: false})
		return
	}

	var captchaResult string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 4; i++ {
		captchaResult += strconv.Itoa(r.Intn(10))
	}

	session.Set(constant.PHONE_CAPTCHA_LAST, time.Now().UnixNano())
	session.Set(constant.PHONE_CAPTCHA, captchaResult)
	session.Set(constant.PHONE_NUMBER, phone)
	session.Save()
	c.JSON(http.StatusOK, util.JsonResult{Data: captchaResult, Success: true, Msg: "成功"})

}
