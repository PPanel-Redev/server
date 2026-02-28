package paymentreceived

import (
	"context"
	"errors"
	"fmt"

	"github.com/perfect-panel/server/pkg/cache"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ Model = (*customPaymentReceivedModel)(nil)

var (
	cachePaymentReceivedIdPrefix = "cache:payment_received:id:"
)

type (
	Model interface {
		paymentReceivedModel
		customPaymentReceivedLogicModel
	}
	paymentReceivedModel interface {
		Insert(ctx context.Context, data *PaymentReceived, tx ...*gorm.DB) error
		FindOne(ctx context.Context, id int64) (*PaymentReceived, error)
		Update(ctx context.Context, data *PaymentReceived, tx ...*gorm.DB) error
		Delete(ctx context.Context, id int64, tx ...*gorm.DB) error
	}

	customPaymentReceivedLogicModel interface {
		FindListByUserId(ctx context.Context, userId int64) ([]*PaymentReceived, error)
		FindListByUserIdAndType(ctx context.Context, userId int64, receivedType string) ([]*PaymentReceived, error)
		QueryListByPage(ctx context.Context, page, size int, userId int64, receivedType string) (int64, []*PaymentReceived, error)
	}

	customPaymentReceivedModel struct {
		*defaultPaymentReceivedModel
	}
	defaultPaymentReceivedModel struct {
		cache.CachedConn
		table string
	}
)

func NewModel(db *gorm.DB, c *redis.Client) Model {
	return &customPaymentReceivedModel{
		defaultPaymentReceivedModel: newPaymentReceivedModel(db, c),
	}
}

func newPaymentReceivedModel(db *gorm.DB, c *redis.Client) *defaultPaymentReceivedModel {
	return &defaultPaymentReceivedModel{
		CachedConn: cache.NewConn(db, c),
		table:      "`payment_received`",
	}
}

func (m *defaultPaymentReceivedModel) getCacheKeys(data *PaymentReceived) []string {
	if data == nil {
		return []string{}
	}
	paymentReceivedIdKey := fmt.Sprintf("%s%v", cachePaymentReceivedIdPrefix, data.Id)
	return []string{paymentReceivedIdKey}
}

func (m *defaultPaymentReceivedModel) Insert(ctx context.Context, data *PaymentReceived, tx ...*gorm.DB) error {
	err := m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		return conn.Create(&data).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *defaultPaymentReceivedModel) FindOne(ctx context.Context, id int64) (*PaymentReceived, error) {
	paymentReceivedIdKey := fmt.Sprintf("%s%v", cachePaymentReceivedIdPrefix, id)
	var resp PaymentReceived
	err := m.QueryCtx(ctx, &resp, paymentReceivedIdKey, func(conn *gorm.DB, v interface{}) error {
		return conn.First(&resp, id).Error
	})
	switch {
	case err == nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *defaultPaymentReceivedModel) Update(ctx context.Context, data *PaymentReceived, tx ...*gorm.DB) error {
	old, err := m.FindOne(ctx, data.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		return conn.Save(data).Error
	}, m.getCacheKeys(old)...)
	return err
}

func (m *defaultPaymentReceivedModel) Delete(ctx context.Context, id int64, tx ...*gorm.DB) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	err = m.ExecCtx(ctx, func(conn *gorm.DB) error {
		if len(tx) > 0 {
			conn = tx[0]
		}
		return conn.Unscoped().Delete(&PaymentReceived{}, id).Error
	}, m.getCacheKeys(data)...)
	return err
}

func (m *customPaymentReceivedModel) FindListByUserId(ctx context.Context, userId int64) ([]*PaymentReceived, error) {
	var resp []*PaymentReceived
	err := m.QueryNoCacheCtx(ctx, &resp, func(conn *gorm.DB, v interface{}) error {
		return conn.Where("`user_id` = ?", userId).Find(v).Error
	})
	return resp, err
}

func (m *customPaymentReceivedModel) FindListByUserIdAndType(ctx context.Context, userId int64, receivedType string) ([]*PaymentReceived, error) {
	var resp []*PaymentReceived
	err := m.QueryNoCacheCtx(ctx, &resp, func(conn *gorm.DB, v interface{}) error {
		db := conn.Where("`user_id` = ?", userId)
		if receivedType != "" {
			db = db.Where("`received_type` = ?", receivedType)
		}
		return db.Find(v).Error
	})
	return resp, err
}

func (m *customPaymentReceivedModel) QueryListByPage(ctx context.Context, page, size int, userId int64, receivedType string) (int64, []*PaymentReceived, error) {
	var list []*PaymentReceived
	var total int64
	err := m.QueryNoCacheCtx(ctx, &list, func(conn *gorm.DB, v interface{}) error {
		conn = conn.Model(&PaymentReceived{})
		if userId > 0 {
			conn = conn.Where("user_id = ?", userId)
		}
		if receivedType != "" {
			conn = conn.Where("received_type = ?", receivedType)
		}
		return conn.Order("id desc").Count(&total).Offset((page - 1) * size).Limit(size).Find(v).Error
	})
	return total, list, err
}
