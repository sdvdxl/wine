package bean

import "time"

// 订单的销售注释
type SalesComments struct {
	Uuid       string
	TimeCreate time.Time
	TimeUpdate time.Time
	OrderUuid  string
	Comments   string
	UserUuid   string
}

type UserSalesComments struct {
	User          `xorm:"extends"`
	SalesComments `xorm:"extends"`
}
