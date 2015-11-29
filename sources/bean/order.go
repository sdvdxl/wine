package bean

import "time"

type OrderInfo struct {
	Uuid                         string
	TimeCreate                   time.Time
	TimeUpdate                   time.Time
	OrderNo                      string    //订单编号
	BuyerName                    string    //买家会员名
	BuyerPayAccount              string    //买家支付宝账号
	BuyerNeedPay                 float64   `json:",string"` //买家应付货款
	BuyerPostage                 float64   //买家应付邮费
	BuyerPayIntegral             int       `json:",string"` //买家支付积分
	TotalPrice                   float64   `json:",string"` //总金额
	ReturnIntegral               int       `json:",string"` // 返点积分
	RealPayAmount                float64   `json:",string"` // 买家实际支付金额
	RealPayIntegral              int       `json:",string"` // 买家实际支付积分
	OrderStatus                  string    // 订单状态
	BuyerMessage                 string    // 买家留言
	DeliveryName                 string    // 收货人姓名
	DeliveryAddressProvince      string    // 收货省份
	DeliveryAddressCity          string    // 收货市
	DeliveryAddressArea          string    // 收货区
	DeliveryAddressTown          string    // 收货镇
	DeliveryAddressStreet        string    // 收货街道
	DeliveryAddressDetail        string    // 详细地址
	DeliveryAddressPost          string    // 邮编
	TransferType                 string    // 运送方式
	ContactTel                   string    // 联系电话
	ContactPhone                 string    // 联系手机
	OrderCreateTime              time.Time // 订单创建时间
	OrderPayTime                 time.Time // 订单付款时间
	GoodsName                    string    // 宝贝标题
	GoodsType                    string    // 宝贝种类
	LogisticsNumber              string    // 物流单号
	LogisticsCompany             string    // 物流公司
	OrderComments                string    // 订单备注
	GoodsAmount                  int       `json:",string"` // 宝贝数量
	ShopId                       int       `json:",string"` // 店铺ID
	ShopName                     string    // 店铺名称
	OrderCloseReason             string    // 订单关闭原因
	SellerServiceFee             float64   `json:",string"` // 卖家服务费
	BuyerServiceFee              float64   `json:",string"` // 买家服务费'
	InvoiceTitle                 string    // 发票抬头
	OrderSource                  string    // 订单来源
	OrderStageInfo               string    // 分阶段订单信息
	FixedDeliveryAddressProvince string    // 收货省份
	FixedDeliveryAddressCity     string    // 收货市
	FixedDeliveryAddressArea     string    // 收货区
	FixedDeliveryAddressTown     string    // 收货镇
	FixedDeliveryAddressStreet   string    // 收货街道
	FixedDeliveryAddressDetail   string    // 详细地址
	FixedDeliveryAddressPost     string    // 邮编

}
