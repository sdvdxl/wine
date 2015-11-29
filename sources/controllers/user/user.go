package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/go-tools/number"
	"github.com/sdvdxl/go-uuid/uuid"
	"github.com/sdvdxl/wine/sources/bean"
	"github.com/sdvdxl/wine/sources/controllers/userauth"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/wine/sources/util/db"
	"github.com/sdvdxl/wine/sources/util/pagehelper"
	"github.com/sdvdxl/wine/sources/util/render"
	"strings"
	"time"
)

const pageSize = 10

// 用户列表
func UserListHandler(c *gin.Context) {
	auth := userauth.Auth(c)
	curUser := auth.CurUser()
	if curUser.Role != bean.ROLE_SUPER_ADMIN {
		data := make(map[string]interface{}, 0)
		data["page"] = pagehelper.Paging(1, 1, 0)
		data["userList"] = make([]bean.User, 0, 0)
		render.New(c).JSON(util.JsonResult{Success: true, Msg: "success", Data: data})
		return
	}

	curPage := number.DefaultInt(c.Query("curPage"), 1)

	query := db.Engine.Where("1=1")
	queryCount := db.Engine.Where("1=1")
	defer query.Close()
	defer queryCount.Close()

	user := new(bean.User)
	count, err := queryCount.Count(user)
	util.PanicError(err)

	query.Limit(pageSize, pageSize*(curPage-1))
	rows, err := query.Cols("uuid", "phone", "nickname", "email", "role", "time_create", "time_update").Rows(user)
	util.PanicError(err)
	defer rows.Close()
	userList := make([]bean.User, 0, 100)
	for rows.Next() {
		user = new(bean.User)
		err = rows.Scan(user)
		userList = append(userList, *user)
	}

	data := make(map[string]interface{})
	data["page"] = pagehelper.Paging(curPage, pageSize, int(count))
	data["userList"] = userList
	render.New(c).JSON(util.JsonResult{Success: true, Msg: "success", Data: data})
}

// 添加用户
func UserAddHandler(c *gin.Context) {
	r := render.New(c)

	//超级管理员才可以添加用户
	auth := userauth.Auth(c)
	if !auth.IsRole(bean.ROLE_SUPER_ADMIN) {
		r.JSON(util.JsonResult{Msg: "您没有权限"})
		return
	}

	//先根据手机号,查询是否已经存在该用户
	phone := strings.TrimSpace(c.PostForm("phone"))
	if len(phone) != 11 {
		r.JSON(util.JsonResult{Msg: "请填写手机号"})
		return
	}

	user := &bean.User{Phone: phone}
	found, err := db.Engine.Get(user)
	util.PanicError(err)
	if found {
		r.JSON(util.JsonResult{Msg: "该手机号已经存在"})
		return
	}

	//校验密码
	user.Phone = phone
	password := c.PostForm("password")
	if len(password) < 6 {
		r.JSON(util.JsonResult{Msg: "密码不能小于6位"})
		return
	}

	//校验角色
	role := strings.TrimSpace(c.PostForm("role"))
	if role == bean.ROLE_ADMIN || role == bean.ROLE_SALES || role == bean.ROLE_SUPER_ADMIN {
		user.Role = role
	} else {
		r.JSON(util.JsonResult{Msg: "请选择用户角色"})
		return
	}

	user.Salt = uuid.New()
	user.Password = util.HashAndSalt(password, user.Salt)

	user.Nickname = c.PostForm("nickname")
	user.Uuid = uuid.New()
	user.TimeCreate = time.Now()
	user.TimeUpdate = time.Now()
	count, err := db.Engine.Insert(user)
	util.PanicError(err)

	if count == 0 {
		r.JSON(util.JsonResult{Msg: "新增用户失败"})
	} else {
		r.JSON(util.JsonResult{Success: true, Msg: "新增用户成功"})
	}
}

// 删除用户
func UserDeleteHandler(c *gin.Context) {
	r := render.New(c)
	//超级管理员才可以删除用户
	auth := userauth.Auth(c)
	if !auth.IsRole(bean.ROLE_SUPER_ADMIN) {
		r.JSON(util.JsonResult{Msg: "您没有权限"})
		return
	}

	userUuid := strings.TrimSpace(c.PostForm("uuid"))
	if userUuid == "" {
		r.JSON(util.JsonResult{Msg: "请选择用户"})
		return
	}

	user := &bean.User{Uuid: userUuid}
	found, err := db.Engine.Get(user)
	util.PanicError(err)
	if !found {
		r.JSON(util.JsonResult{Msg: "没有此用户"})
		return
	}

	//不可以删除自己
	if auth.CurUser().Uuid == user.Uuid {
		r.JSON(util.JsonResult{Msg: "不可以删除自己"})
		return
	}

	user = &bean.User{Uuid: userUuid}
	count, err := db.Engine.Delete(user)
	util.PanicError(err)
	if count == 0 {
		r.JSON(util.JsonResult{Msg: "删除失败"})
		return
	}

	r.JSON(util.JsonResult{Success: true, Msg: "删除成功"})
}

//更改用户
func UserUpdateHandler(c *gin.Context) {
	r := render.New(c)

	//超级管理员才可以修改用户
	auth := userauth.Auth(c)
	if !auth.IsRole(bean.ROLE_SUPER_ADMIN) {
		r.JSON(util.JsonResult{Msg: "您没有权限"})
		return
	}

	uuid := strings.TrimSpace(c.PostForm("uuid"))
	if uuid == "" {
		r.JSON(util.JsonResult{Msg: "请选择用户"})
		return
	}

	user := &bean.User{Uuid: uuid}
	found, err := db.Engine.Get(user)
	util.PanicError(err)
	if !found {
		r.JSON(util.JsonResult{Msg: "用户不存在"})
		return
	}

	//如果是管理员自己不允许更改角色,防止超级管理员身份失效
	isSuperAdminSelf := false
	if uuid == auth.CurUser().Uuid {
		isSuperAdminSelf = true
	}

	phone := strings.TrimSpace(c.PostForm("phone"))
	user = &bean.User{Phone: phone}
	found, err = db.Engine.Get(user)
	util.PanicError(err)
	if found && user.Uuid != uuid {
		r.JSON(util.JsonResult{Msg: "手机号已经存在"})
		return
	}

	userCond := &bean.User{Uuid: uuid}
	user = &bean.User{Phone: phone, Password: c.PostForm("password"), Nickname: strings.TrimSpace(c.PostForm("nickname"))}
	if !isSuperAdminSelf {
		user.Role = c.PostForm("role")
	}

	_, err = db.Engine.Update(user, userCond)
	util.PanicError(err)

	r.JSON(util.JsonResult{Success: true, Msg: "修改成功"})
}

// 我的信息
func MyProfileHandler(c *gin.Context) {
	r := render.New(c)
	auth := userauth.Auth(c)
	curUser := auth.CurUser()
	user := bean.User{Phone: curUser.Phone, Nickname: curUser.Nickname, Role: curUser.Role}
	r.JSON(util.JsonResult{Success: true, Msg: "获取成功", Data: user})
}

//更改密码
func ChangePasswordHandler(c *gin.Context) {
	r := render.New(c)
	oldPassword := c.PostForm("old")
	if oldPassword == "" {
		r.JSON(util.JsonResult{Msg: "请输入旧密码"})
		return
	}

	auth := userauth.Auth(c)
	curUser := auth.CurUser()
	if curUser.Password != util.HashAndSalt(oldPassword, curUser.Salt) {
		r.JSON(util.JsonResult{Msg: "原密码不正确"})
		return
	}

	newPassword := c.PostForm("new")
	if len(newPassword) < 6 {
		r.JSON(util.JsonResult{Msg: "新密码长度不能小于6位"})
		return
	}

	userCond := &bean.User{Uuid: auth.CurUser().Uuid}
	user := &bean.User{Salt: uuid.New(), TimeUpdate: time.Now()}
	user.Password = util.HashAndSalt(newPassword, user.Salt)
	count, err := db.Engine.Update(user, userCond)
	util.PanicError(err)

	if count == 0 {
		r.JSON(util.JsonResult{Msg: "更改失败"})
		return
	}

	auth.Logout()
	r.JSON(util.JsonResult{Success: true, Msg: "更改成功"})

}
