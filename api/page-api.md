# 获取用户收款方式列表（分页）

## 接口信息

- **接口路径**: `/v1/admin/paymentReceived/page`
- **请求方式**: `GET`
- **接口描述**: 管理员分页查询用户的收款方式列表，支持按用户ID和收款类型筛选

## 请求参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int64 | 是 | 页码，从1开始 |
| size | int64 | 是 | 每页数量 |
| user_id | int64 | 否 | 用户ID，不传则查询所有用户 |
| received_type | string | 否 | 收款类型，可选值：`weixin`(微信)、`alipay`(支付宝)、`bankcard`(银行卡) |

## 请求示例

```bash
curl -X GET "https://api.example.com/v1/admin/paymentReceived/page?page=1&size=10&user_id=123&received_type=weixin" \
  -H "Authorization: Bearer <token>"
```

## 响应参数

| 参数名 | 类型 | 说明 |
|--------|------|------|
| total | int64 | 总记录数 |
| list | array | 收款方式列表 |

### list 数组元素

| 参数名 | 类型 | 说明 |
|--------|------|------|
| id | int64 | 收款方式ID |
| user_id | int64 | 用户ID |
| received_name | string | 收款人姓名 |
| received_type | string | 收款类型：`weixin`(微信)、`alipay`(支付宝)、`bankcard`(银行卡) |
| received_no | string | 收款账号/银行账号/卡号 |
| bank_name | string | 银行名称（银行卡类型时有值） |
| opening_branch | string | 开户行（银行卡类型时有值） |
| qrcode | string | 收款二维码链接（微信/支付宝类型时有值） |
| created_at | int64 | 创建时间（Unix时间戳，秒） |
| updated_at | int64 | 更新时间（Unix时间戳，秒） |

## 响应示例

### 成功响应

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "total": 2,
    "list": [
      {
        "id": 1,
        "user_id": 123,
        "received_name": "张三",
        "received_type": "weixin",
        "received_no": "wx_abc123",
        "bank_name": "",
        "opening_branch": "",
        "qrcode": "https://example.com/qrcode/weixin_1.png",
        "created_at": 1708900000,
        "updated_at": 1708900000
      },
      {
        "id": 2,
        "user_id": 123,
        "received_name": "张三",
        "received_type": "bankcard",
        "received_no": "6222021234567890123",
        "bank_name": "中国工商银行",
        "opening_branch": "北京分行",
        "qrcode": "",
        "created_at": 1708900000,
        "updated_at": 1708900000
      }
    ]
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "msg": "参数错误",
  "data": null
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权/Token无效 |
| 500 | 服务器内部错误 |

## 注意事项

1. 该接口需要管理员权限，需要通过 `Authorization` 请求头传递有效的 JWT Token
2. `received_type` 参数不传时返回所有类型的收款方式
3. `user_id` 参数不传时返回所有用户的收款方式
4. 时间字段为 Unix 时间戳（秒级）