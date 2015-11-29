CREATE TABLE `user` (
  `uuid` char(36) PRIMARY KEY  COMMENT 'uuid',
  `time_create` datetime NOT NULL COMMENT '创建时间',
  `time_update` datetime NOT NULL COMMENT '更新时间',
  `nickname` varchar(64) DEFAULT NULL COMMENT '用户昵称',
  `password` varchar(512) NOT NULL COMMENT '密码',
  `salt` varchar(512) NOT NULL COMMENT '盐',
  `phone` varchar(13) unique NOT NULL COMMENT '手机号',
  `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
  `role` varchar(64) DEFAULT 'sales' COMMENT '角色,默认是销售,sales,admin super_admin'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户';

create table `order_info` (
  `uuid` char(36) primary key,
  `time_create` datetime not null comment '创建时间',
  `time_update` datetime not null comment '修改时间',
  `order_no` varchar(100) comment '订单编号',
  `buyer_name` varchar(64) comment '买家会员名',
  `buyer_pay_account` varchar(64) comment '买家支付宝账号',
  `buyer_need_pay` decimal(10,2) default 0 comment '买家应付货款',
  `buyer_postage` decimal(10,2) default 0 comment '买家应付邮费',
  `buyer_pay_integral` int default 0 comment '买家支付积分',
  `total_price` decimal(10,2) default 0 comment '总金额',
  `return_integral` int default 0 comment '返点积分',
  `real_pay_amount` decimal(10,2) default 0 comment '买家实际支付金额',
  `real_pay_integral` int default 0 comment '买家实际支付积分',
  `order_status` varchar(200) comment '订单状态, 等待买家付款，买家已付款，等待卖家发货，卖家已发货，等待买家确认，交易关闭，交易成功',
  `buyer_message` varchar(2000) default '' comment '买家留言',
  `delivery_name` varchar(64) comment '收货人姓名',
  `delivery_address_province` varchar(1000)  comment '收货省份',
  `delivery_address_city` varchar(1000) comment '收货市',
  `delivery_address_area` varchar(1000) comment '收货区',
  `delivery_address_town` varchar(1000)  comment '收货镇',
  `delivery_address_street` varchar(1000) comment '收货街道',
  `delivery_address_detail` varchar(1000) comment '详细地址',
  `delivery_address_post` varchar(6) comment '邮编',
  `transfer_type` varchar(100) comment '运送方式',
  `contact_tel` varchar(20) comment '联系电话',
  `contact_phone` varchar(20) comment '联系手机',
  `order_create_time` datetime comment '订单创建时间',
  `order_pay_time` datetime comment '订单付款时间',
  `goods_name` varchar(1000) comment '宝贝标题',
  `goods_type` varchar(10) comment '宝贝种类',
  `logistics_number` varchar(64) comment '物流单号',
  `logistics_company` varchar(200) comment '物流公司',
  `order_comments` varchar(2000) comment '订单备注',
  `goods_amount` int default 0 comment '宝贝数量',
  `shop_id` int default 0 comment '店铺ID',
  `shop_name` varchar(200) comment '店铺名称',
  `order_close_reason` varchar(200) comment '订单关闭原因',
  `seller_service_fee` decimal(10,2) default 0 comment '卖家服务费',
  `buyer_service_fee` decimal(10,2) default 0 comment '买家服务费',
  `invoice_title` varchar(200) comment '发票抬头',
  `order_source` varchar(100) comment '订单来源',
  `order_stage_info` varchar(500) comment '分阶段订单信息',
  `fixed_delivery_address_province` varchar(1000) comment '修改后的收货省份',
  `fixed_delivery_address_city` varchar(1000) comment '修改后的收货市',
  `fixed_delivery_address_area` varchar(1000) comment '修改后的收货区',
  `fixed_delivery_address_town` varchar(1000) comment '修改后的收货镇',
  `fixed_delivery_address_street` varchar(1000) comment '修改后的收货街道',
  `fixed_delivery_address_detail` varchar(1000) comment '修改后的详细地址',
  `fixed_delivery_address_post` varchar(6) comment '邮编',
  unique key `idx_order_no` (`order_no`)
)


create table sales_comments (
	uuid char(36) primary key comment '主键',
	time_create datetime not null comment '创建时间',
	time_update datetime not null comment '修改时间',
	order_uuid char(36) not null comment '订单的主键',
	user_uuid char(36) not null comment '添加用户主键' ,
	comments varchar(2000) not null comment '销售注释'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='销售注释';

------ 初始化语句
insert into `user`(`uuid`, `time_create`, `time_update`, `nickname`, `password`,`salt`, `phone`, `role` )
values('ddd4b5a9-fecd-446c-bd78-63b70bb500a1', now(), now(), '黄总',
'1664ad4059146142c15288688b4e138b066b27463fadad8009ab038aafa6f5016060f0360f4aab418e96de5b5f3e43349d9d6897d372b69a6a71a72b3994cdfc',
'245accec-3c12-4642-967f-e476cef558c4', '13809777237', 'super_admin');

insert into `user`(`uuid`, `time_create`, `time_update`, `nickname`, `password`,`salt`, `phone`, `role` )
values('edd4b5a9-fecd-446c-bd78-63b70bb500a1', now(), now(), '杜龙少',
'cea4ca9d951bedc59f1cad5b6c134202843821ad0e549625efb9c25767417b4fa044a19a47220fd1c227d2dc4cbaefbdc86f4cdef0845e45cb618b2e1cb6f5a3',
'245accec-3c12-4642-967f-e476cef558c0', '13165360918', 'super_admin');