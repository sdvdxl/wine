package order

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sdvdxl/go-tools/number"
	timeutils "github.com/sdvdxl/go-tools/time"
	"github.com/sdvdxl/go-uuid/uuid"
	"github.com/sdvdxl/wine/sources/bean"
	"github.com/sdvdxl/wine/sources/controllers/userauth"
	"github.com/sdvdxl/wine/sources/util"
	"github.com/sdvdxl/wine/sources/util/db"
	. "github.com/sdvdxl/wine/sources/util/log"
	"github.com/sdvdxl/wine/sources/util/pagehelper"
	"github.com/sdvdxl/wine/sources/util/render"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const pageSize = 10

func OrderListHandler(c *gin.Context) {
	r := render.New(c)
	orderInfo := new(bean.OrderInfo)

	curPage := number.DefaultInt(c.Query("curPage"), 1)
	seaType := strings.TrimSpace(c.Query("seaType"))
	seaContent := strings.TrimSpace(c.Query("seaContent"))

	query := db.Engine.Where("1=1")
	queryCount := db.Engine.Where("1=1")
	defer query.Close()
	defer queryCount.Close()

	//订单时间
	orderStart := timeutils.ParseDate(c.Query("orderStart"))
	Logger.Debug("order start time: %v", orderStart)
	if orderStart != nil {
		queryCount.And("order_create_time>=?", orderStart)
		query.And("order_create_time>=?", orderStart)
	}

	orderEnd := timeutils.ParseDate(c.Query("orderEnd"))
	Logger.Debug("order end time: %v", orderEnd)
	if orderEnd != nil {
		t := orderEnd.Add(time.Hour * 24)
		orderEnd = &t
		queryCount.And("order_create_time<?", orderEnd)
		query.And("order_create_time<?", orderEnd)
	}

	//搜索类型和内容
	if seaType == "date-baught-lt" || seaType == "date-baught-gt" || seaType == "date-baught-eq" {
		orderCount := number.DefaultInt(seaContent, -1)
		if -1 == orderCount {
			r.JSON(util.JsonResult{Msg: "请输入次数(整数)"})
			return
		}

		//购买次数少于
		operator := ">"
		if seaType == "date-baught-lt" {
			operator = "<"
		} else if seaType == "date-baught-eq" {
			operator = "="
		}

		query.Sql(`select * from (select * from order_info where buyer_name in  (select buyer_name from
(select *, count(buyer_name) c from order_info group by buyer_name having c `+operator+`?) tmp_table) ) tmp_count_table order by order_create_time desc limit ?, ?`, orderCount, pageSize*(curPage-1), pageSize)
		queryCount.Sql(`select count(*) from (select * from order_info where buyer_name in  (select buyer_name from
(select *, count(buyer_name) c from order_info group by buyer_name having c `+operator+`?) tmp_table) ) tmp_count_table `, orderCount)
	} else if seaType != "" && seaContent != "" {
		seaContent = "%" + seaContent + "%"
		if "deliveryAddress" == seaType {
			query.And(`delivery_address_province like ? or delivery_address_city like ?  or delivery_address_area like ?
		 	 or delivery_address_town like ?  or delivery_address_street like ? or delivery_address_detail like ? or delivery_address_post like ? `,
				seaContent, seaContent, seaContent, seaContent, seaContent, seaContent, seaContent).OrderBy("order_create_time desc")
		} else {
			switch seaType {
			case "buyerName":
				seaType = "buyer_name"
			case "deliveryName":
				seaType = "delivery_name"
			case "buyerPayAccount":
				seaType = "buyer_pay_account"
			case "contactTel":
				seaType = "contact_tel"
			case "contactPhone":
				seaType = "contact_phone"
			}
			query.And(seaType+" like ? ", seaContent)
			queryCount.And(seaType+" like ? ", seaContent)
		}
	} else if seaContent != "" {
		whereCond := make([]interface{}, 22, 22)
		for i, _ := range whereCond {
			whereCond[i] = "%" + seaContent + "%"
		}

		queryCondSQL := `buyer_name like ? or buyer_pay_account like ? or  buyer_need_pay like ? or
		buyer_postage like ? or  buyer_pay_integral like ? or  total_price like ? or 	return_integral like ? or
		real_pay_amount like ? or  real_pay_integral like ? or  order_status like ? or
		buyer_message like ? or  delivery_name like ? or  delivery_address_province like ? or
		delivery_address_city like ? or  delivery_address_area like ? or  delivery_address_town like ? or
		delivery_address_street like ? or  delivery_address_detail like ? or
		delivery_address_post like ? or  transfer_type like ? or  contact_tel like ? or  contact_phone like ?`

		queryCount.Where(queryCondSQL, whereCond...)
		query.Where(queryCondSQL, whereCond...)
	}

	count, err := queryCount.Count(orderInfo)
	util.PanicError(err)

	query.Limit(pageSize, pageSize*(curPage-1))
	rows, err := query.Desc("order_create_time").Rows(orderInfo)
	util.PanicError(err)

	defer rows.Close()
	orderList := make([]bean.OrderInfo, 0, 100)

	for rows.Next() {
		orderInfo = new(bean.OrderInfo)
		err = rows.Scan(orderInfo)
		orderList = append(orderList, *orderInfo)
	}

	data := make(map[string]interface{})
	data["page"] = pagehelper.Paging(curPage, pageSize, int(count))
	data["orderList"] = orderList
	r.JSON(util.JsonResult{Success: true, Msg: "success", Data: data})
}

// 更新订单
func OrderUpdateHandler(c *gin.Context) {
	r := render.New(c)

	//超级管理员才可以添加用户
	auth := userauth.Auth(c)
	if !auth.IsRole(bean.ROLE_SUPER_ADMIN) {
		r.JSON(util.JsonResult{Msg: "您没有权限"})
		return
	}

	var orderInfo = new(bean.OrderInfo)
	err := json.Unmarshal([]byte(c.PostForm("orderInfo")), orderInfo)
	util.PanicError(err)
	Logger.Debug("bind order : %v", orderInfo)

	if orderInfo.Uuid == "" {
		r.JSON(util.JsonResult{Msg: "请选择订单"})
		return
	}

	if orderInfo.OrderNo == "" {
		r.JSON(util.JsonResult{Msg: "请输入订单编号"})
		return
	}

	checkOrder := &bean.OrderInfo{OrderNo: orderInfo.OrderNo}
	found, err := db.Engine.Get(checkOrder)
	util.PanicError(err)
	if found && checkOrder.Uuid != orderInfo.Uuid {
		r.JSON(util.JsonResult{Msg: "已经存在相同的订单编号"})
		return
	}

	orderInfo.TimeUpdate = time.Now()
	address := strings.Split(orderInfo.DeliveryAddressProvince, " ")
	if len(address) == 4 {
		orderInfo.DeliveryAddressProvince = address[0]
		orderInfo.DeliveryAddressCity = address[1]
		orderInfo.DeliveryAddressArea = address[2]
		orderInfo.DeliveryAddressStreet = address[3]
	} else {
		orderInfo.DeliveryAddressProvince = ""
	}
	conditionBean := bean.OrderInfo{Uuid: orderInfo.Uuid}
	_, err = db.Engine.Update(orderInfo, conditionBean)
	util.PanicError(err)

	r.JSON(util.JsonResult{Msg: "更新成功", Success: true})
}

// 追加备注
func OrderAddCommentsHandler(c *gin.Context) {
	r := render.New(c)
	oid := strings.TrimSpace(c.PostForm("oid"))
	if oid == "" {
		r.JSON(util.JsonResult{Msg: "请选择订单"})
		return
	}

	order := &bean.OrderInfo{Uuid: oid}
	found, err := db.Engine.Get(order)
	util.PanicError(err)
	if !found {
		r.JSON(util.JsonResult{Msg: "没有此订单"})
		return
	}

	comments := strings.TrimSpace(c.PostForm("comments"))
	if comments == "" {
		r.JSON(util.JsonResult{Msg: "请填写内容"})
		return
	}

	auth := userauth.Auth(c)
	salesComments := &bean.SalesComments{Uuid: uuid.New(),
		TimeCreate: time.Now(),
		TimeUpdate: time.Now(),
		Comments:   comments,
		UserUuid:   auth.CurUser().Uuid,
		OrderUuid:  oid,
	}

	count, err := db.Engine.Insert(salesComments)
	util.PanicError(err)
	if count == 0 {
		r.JSON(util.JsonResult{Msg: "添加注释失败"})
		return
	}

	r.JSON(util.JsonResult{Success: true, Msg: "添加成功"})
}

// 根据订单uuid查看注释
func OrderCommentsHandler(c *gin.Context) {
	r := render.New(c)
	oid := strings.TrimSpace(c.Query("oid"))
	if oid == "" {
		r.JSON(util.JsonResult{Msg: "请选择订单"})
		return
	}

	order := &bean.OrderInfo{Uuid: oid}
	found, err := db.Engine.Get(order)
	util.PanicError(err)
	if !found {
		r.JSON("没有此订单")
		return
	}

	comments := make([]bean.UserSalesComments, 0, 30)
	err = db.Engine.Cols("u.uuid", "u.nickname", "u.phone", "comments").Table("user").Alias("u").Join("inner", "sales_comments", "u.uuid=user_uuid").Where("order_uuid=?", oid).Find(&comments)
	util.PanicError(err)

	r.JSON(util.JsonResult{Success: true, Msg: "获取成功", Data: comments})
}

const (
	insertSQL = `insert ignore into order_info(uuid,time_create, time_update, order_no, buyer_name, buyer_pay_account, buyer_need_pay,
buyer_postage, buyer_pay_integral, total_price,	return_integral, real_pay_amount, real_pay_integral, order_status,
buyer_message, delivery_name, delivery_address_province, delivery_address_city, delivery_address_area, delivery_address_town, delivery_address_street, delivery_address_detail,
delivery_address_post, transfer_type, contact_tel, contact_phone, order_create_time, order_pay_time, goods_name, goods_type, logistics_number,
logistics_company, order_comments, goods_amount, shop_id, shop_name, order_close_reason, seller_service_fee, buyer_service_fee, invoice_title,
order_source, order_stage_info, fixed_delivery_address_province, fixed_delivery_address_city, fixed_delivery_address_area, fixed_delivery_address_town,
fixed_delivery_address_street, fixed_delivery_address_detail, fixed_delivery_address_post) values`
)

func OrderUploadTemplateHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "order_upload.tmpl", nil)
}

// 上传订单
func OrderUploadHandler(c *gin.Context) {

	//销售不允许上传
	auth := userauth.Auth(c)
	if auth.IsRole(bean.ROLE_SALES) {
		c.HTML(http.StatusOK, "order_upload.tmpl", gin.H{"successCount": 0, "isUpload": true, "errMsg": "您没有权限"})
		return
	}

	request := c.Request
	request.ParseMultipartForm(2 << 10)
	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		Logger.Error("upload file error, %v", err)
		panic(err)
	}

	_ = fileHeader

	csvReader := csv.NewReader(file)

	sql := insertSQL
	values := make([]interface{}, 0, 40)

	_, err = csvReader.Read()
	if err != nil {
		panic(err)
	} else if err == io.EOF {
		Logger.Debug("read file end")
		return
	}

	rowsCount := 0
	for {
		if row, err := csvReader.Read(); err != nil && err == io.EOF {
			break
		} else if err != io.EOF {
			rowsCount++
			// 49 params  47 ?
			sql += `(?, now(), now(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? ),`
			values = append(values, uuid.New())
			r := new(big.Rat)
			r.SetString(row[0])
			values = append(values, r.RatString())
			values = append(values, appendValues(row, 1, 12)...)
			Logger.Debug("address: %v", row[13])
			address := strings.Fields(row[13]) //dizhi

			Logger.Debug("address args len:%v", len(address))
			if len(address) == 4 {
				values = append(values, address[0])
				values = append(values, address[1])
				values = append(values, address[2])
				values = append(values, address[3])
				for i := 0; i < 3; i++ {
					values = append(values, nil)
				}
			} else {
				for i := 0; i < 7; i++ {
					values = append(values, nil)
				}
			}

			values = append(values, appendValues(row, 14, 15)...)

			runes := []rune(row[16])
			if strings.HasPrefix(row[16], "'") {
				runes = runes[1:]
			}
			values = append(values, string(runes))

			values = append(values, appendValues(row, 17, 28)...)

			runes = []rune(row[29])
			if strings.HasSuffix(row[29], "元") {
				runes = runes[:len(runes)-1]
			}
			buyerServiceFee, err := strconv.Atoi(string(runes))
			if err != nil {
				buyerServiceFee = 0
			}
			values = append(values, buyerServiceFee)
			values = append(values, appendValues(row, 30, 32)...)

			address = strings.Fields(row[35]) //dizhirow[36] //fixed dizhi
			Logger.Debug("fixed address: %v", row[35])
			Logger.Debug("fixed address args len:%v", len(address))
			// 7
			if len(address) == 4 {
				values = append(values, address[0])
				values = append(values, address[1])
				values = append(values, address[2])
				values = append(values, address[3])
				for i := 0; i < 3; i++ {
					values = append(values, nil)
				}
			} else {
				for i := 0; i < 7; i++ {
					values = append(values, nil)
				}
			}
		}
	}

	sql = sql[:len(sql)-1]

	Logger.Debug("sql:%v", sql)
	Logger.Debug("values len:%v", len(values))
	Logger.Debug("values:%v", values)

	if len(values) > 0 {
		stmt, err := db.Engine.DB().Prepare(sql)
		util.PanicError(err)
		defer stmt.Close()
		result, err := stmt.Exec(values...)
		util.PanicError(err)

		count, err := result.RowsAffected()
		util.PanicError(err)
		Logger.Debug("inserted rows count:%v", count)
	}

	c.HTML(http.StatusOK, "order_upload.tmpl", gin.H{"successCount": rowsCount, "isUpload": true, "errMsg": "上传成功"})

}

func appendValues(row []string, start, end int) []interface{} {
	values := make([]interface{}, 0, 10)
	for i := start; i <= end; i++ {
		Logger.Debug("column %v, value:%v", i, row[i])
		if fmt.Sprint(row[i]) == "" {
			values = append(values, nil)
		} else {
			values = append(values, row[i])
		}

	}

	return values
}

// 删除订单
func OrderDeleteHandler(c *gin.Context) {
	r := render.New(c)

	//超级管理员才可以删除订单
	auth := userauth.Auth(c)
	if !auth.IsRole(bean.ROLE_SUPER_ADMIN) {
		r.JSON(util.JsonResult{Msg: "您没有权限"})
		return
	}

	oids := strings.TrimSpace(c.PostForm("oid"))

	var orderIds []string
	err := json.Unmarshal([]byte(oids), &orderIds)
	if err != nil {
		Logger.Warn("order delete uuids:%v, err: %v", oids, err)
		r.JSON(util.JsonResult{Success: true, Msg: "订单参数错误"})
		return
	}

	if len(orderIds) == 0 {
		r.JSON(util.JsonResult{Success: true, Msg: "请选择订单"})
		return
	}

	deleteSQL := "delete from order_info where uuid in ("
	for _, v := range orderIds {
		deleteSQL += "'" + v + "' ,"
	}
	deleteSQL = deleteSQL[:len(deleteSQL)-1] + ")"
	Logger.Debug("order delete sql: %v", deleteSQL)
	_, err = db.Engine.Exec(deleteSQL)
	util.PanicError(err)
	r.JSON(util.JsonResult{Success: true, Msg: "删除成功"})
}
