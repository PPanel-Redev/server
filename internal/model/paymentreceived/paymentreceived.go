package paymentreceived

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type PaymentReceived struct {
	Id            int64                 `gorm:"primaryKey"`
	UserId        int64                 `gorm:"type:bigint;not null;comment:用户id"`
	ReceivedName  string                `gorm:"type:varchar(100);not null;comment:收款人姓名"`
	ReceivedType  string                `gorm:"type:varchar(20);not null;comment:收款方式;weixin-微信;alipay-支付宝;bankcard-银行卡"`
	ReceivedNo    string                `gorm:"type:varchar(200);not null;comment:收款账号/银行账号/卡号"`
	BankName      string                `gorm:"type:varchar(200);comment:银行名称"`
	OpeningBranch string                `gorm:"type:varchar(200);comment:开户支行"`
	Qrcode        string                `gorm:"type:longtext;comment:二维码;base64格式"`
	CreatedAt     time.Time             `gorm:"<-:create;comment:创建时间"`
	UpdatedAt     time.Time             `gorm:"comment:更新时间"`
	DeletedAt     gorm.DeletedAt        `gorm:"index;comment:删除时间"`
	IsDel         soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt;comment:1:正常 0:删除"`
}

func (*PaymentReceived) TableName() string {
	return "payment_received"
}

const (
	ReceivedTypeWeixin   = "weixin"
	ReceivedTypeAlipay   = "alipay"
	ReceivedTypeBankcard = "bankcard"
)
