# 收款方式

## 表设计
```sql
CREATE TABLE payment_received(
    id bigint AUTO_INCREMENT PRIMARY KEY,
    user_id bigint NOT NULL COMMENT '用户id;user#id',
    received_type varchar(20) NOT NULL COMMENT '收款方式;weixin-微信;alipay-支付宝;-bankcard-银行卡',
    received_no varchar(200) NOT NULL COMMENT '收款账号/银行账号/卡号',
    bank_name varchar(200) COMMENT '银行名称',
    opening_branch varchar(200) COMMENT '开户支行',
    qrcode longtext COMMENT '二维码;base64格式',
    created_at datetime NOT NULL COMMENT '创建时间',
    updated_at datetime COMMENT '修改时间',
    del_at datetime COMMENT '删除时间',
    is_del bigint unsigned NOT NULL COMMENT '1: Normal 0: Deleted'
) COMMENT = '收款方式'
```

## 任务流程
1. 阅读`/Users/weicong/data/code/go/server/internal/handler/admin/order/createOrderHandler.go`以及调用的logic了解handler和logic书写风格
2. 在`/Users/weicong/data/code/go/server/internal/types/types.go`下编写相关请求类、响应类和实体类, **注意这个文件很长，所以你应该只读取最后100行代码了解书写风格即可**
3. 阅读本章节 `接口开发` 
4. 在`/Users/weicong/data/code/go/server/internal/logic/admin/paymentreceived`目录下编写logic相关的代码，一个接口一个文件
5. 在`/Users/weicong/data/code/go/server/internal/handler/admin/paymentreceived`目录下编写handler相关的代码

## 注意事项
 - 每个接口你应该开启subagent来进行这样会让你的速度快很多

## 接口开发
### 1. 添加收款方式
对`payment_received`表进行添加，需要的注意事项如下：
**添加银行卡**
 - user_id: 当前登录人的user_id
 - received_no: 为必填项
 - bank_name: 为必填项
 - opening_branch: 为选填项
 - qrcode: 不需要填写

**添加微信或者支付宝**
- user_id: 当前登录人的user_id
- received_no: 为必填项
- bank_name: 不需要填写
- opening_branch: 不需要填写
- qrcode: 为必填项

## 2.修改收款方式
1. 与添加收款方式一致，但是需要记录修改时间
2. 需要判断是否被逻辑删除

## 3.删除收款方式
对收款方式进行逻辑删除

## 4.查询收款方式列表(不分页)
查询当前用户所有的收款方式，需要过滤逻辑删除的收款方式