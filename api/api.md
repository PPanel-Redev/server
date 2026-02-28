# 收款方式管理接口文档

**基础路径**: `/v1/admin/paymentReceived`

**认证方式**: 所有接口需要在 Header 中携带有效的 Authorization Token

---

## 1. 创建收款方式

**POST** `/v1/admin/paymentReceived/`

### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| received_type | string | 是 | 收款方式类型：`weixin` / `alipay` / `bankcard` |
| received_no | string | 是 | 收款账号（微信号/支付宝账号/银行卡号） |
| bank_name | string | 条件必填 | 银行名称（bankcard 类型必填） |
| opening_branch | string | 否 | 开户支行（仅 bankcard 类型有效） |
| qrcode | string | 条件必填 | 收款二维码链接（weixin/alipay 类型必填） |

### 请求示例

**银行卡类型：**
```json
{
    "received_type": "bankcard",
    "received_no": "6222021234567890123",
    "bank_name": "中国工商银行",
    "opening_branch": "北京分行朝阳支行"
}
```

**微信/支付宝类型：**
```json
{
    "received_type": "weixin",
    "received_no": "wxid_abc123",
    "qrcode": "https://example.com/qrcode/weixin.png"
}
```

### 响应

**成功：**
```json
{
    "code": 0,
    "msg": "success"
}
```

**失败：**
```json
{
    "code": 400,
    "msg": "received_no is required for bankcard"
}
```

---

## 2. 删除收款方式

**POST** `/v1/admin/paymentReceived/del`

### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 收款方式 ID |

### 请求示例

```json
{
    "id": 1
}
```

### 响应

**成功：**
```json
{
    "code": 0,
    "msg": "success"
}
```

**失败：**
```json
{
    "code": 400,
    "msg": "payment received not found"
}
```

---

## 3. 更新收款方式

**PUT** `/v1/admin/paymentReceived/`

### 请求参数

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| id | int64 | 是 | 收款方式 ID |
| received_type | string | 是 | 收款方式类型：`weixin` / `alipay` / `bankcard` |
| received_no | string | 是 | 收款账号 |
| bank_name | string | 条件必填 | 银行名称（bankcard 类型必填） |
| opening_branch | string | 否 | 开户支行 |
| qrcode | string | 条件必填 | 收款二维码链接（weixin/alipay 类型必填） |

### 请求示例

```json
{
    "id": 1,
    "received_type": "alipay",
    "received_no": "13800138000",
    "qrcode": "https://example.com/qrcode/alipay.png"
}
```

### 响应

**成功：**
```json
{
    "code": 0,
    "msg": "success"
}
```

**失败：**
```json
{
    "code": 400,
    "msg": "payment received not found"
}
```

---

## 4. 获取收款方式列表

**GET** `/v1/admin/paymentReceived/`

### 请求参数

无（从 Token 中获取当前用户 ID）

### 响应

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "list": [
            {
                "id": 1,
                "user_id": 100,
                "received_type": "bankcard",
                "received_no": "6222021234567890123",
                "bank_name": "中国工商银行",
                "opening_branch": "北京分行朝阳支行",
                "qrcode": "",
                "created_at": 1709000000,
                "updated_at": 1709100000
            },
            {
                "id": 2,
                "user_id": 100,
                "received_type": "weixin",
                "received_no": "wxid_abc123",
                "bank_name": "",
                "opening_branch": "",
                "qrcode": "https://example.com/qrcode/weixin.png",
                "created_at": 1709000000,
                "updated_at": 0
            }
        ]
    }
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| id | int64 | 收款方式 ID |
| user_id | int64 | 用户 ID |
| received_type | string | 收款方式类型 |
| received_no | string | 收款账号 |
| bank_name | string | 银行名称（bankcard 类型有值） |
| opening_branch | string | 开户支行（bankcard 类型可能有值） |
| qrcode | string | 收款二维码链接（weixin/alipay 类型有值） |
| created_at | int64 | 创建时间（Unix 时间戳） |
| updated_at | int64 | 更新时间（Unix 时间戳，未更新时为 0） |

---

## 5. 按类型分组获取收款方式列表

**GET** `/v1/admin/paymentReceived/listByType`

### 请求参数

无（从 Token 中获取当前用户 ID）

### 响应

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "received_type": [
            {
                "type": "bankcard",
                "list": [
                    {
                        "id": 1,
                        "user_id": 100,
                        "received_type": "bankcard",
                        "received_no": "6222021234567890123",
                        "bank_name": "中国工商银行",
                        "opening_branch": "北京分行朝阳支行",
                        "qrcode": "",
                        "created_at": 1709000000,
                        "updated_at": 0
                    }
                ]
            },
            {
                "type": "weixin",
                "list": [
                    {
                        "id": 2,
                        "user_id": 100,
                        "received_type": "weixin",
                        "received_no": "wxid_abc123",
                        "bank_name": "",
                        "opening_branch": "",
                        "qrcode": "https://example.com/qrcode/weixin.png",
                        "created_at": 1709000000,
                        "updated_at": 0
                    }
                ]
            }
        ]
    }
}
```

### 响应字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| received_type | array | 按类型分组的收款方式数组 |
| received_type[].type | string | 收款方式类型 |
| received_type[].list | array | 该类型下的收款方式列表 |

---

## 字段填写规则

| received_type | received_no | bank_name | opening_branch | qrcode |
|---------------|-------------|-----------|----------------|--------|
| `weixin` | 必填 | 不填 | 不填 | 必填 |
| `alipay` | 必填 | 不填 | 不填 | 必填 |
| `bankcard` | 必填 | 必填 | 可选 | 不填 |

---

## 错误码说明

| code | 说明 |
|------|------|
| 0 | 成功 |
| 400 | 参数错误（缺少必填字段、类型不正确等） |
| 401 | 未授权（Token 无效或过期） |
| 500 | 服务器内部错误 |