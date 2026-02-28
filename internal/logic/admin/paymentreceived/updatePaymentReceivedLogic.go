package paymentreceived

import (
	"context"
	"time"

	"github.com/perfect-panel/server/internal/model/paymentreceived"
	"github.com/perfect-panel/server/internal/model/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/constant"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdatePaymentReceivedLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePaymentReceivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaymentReceivedLogic {
	return &UpdatePaymentReceivedLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePaymentReceivedLogic) UpdatePaymentReceived(req *types.UpdatePaymentReceivedRequest) error {
	u, ok := l.ctx.Value(constant.CtxKeyUser).(*user.User)
	if !ok {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	userId := u.Id
	// 查询收款方式是否存在
	data, err := l.svcCtx.PaymentReceivedModel.FindOne(l.ctx, req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "payment received not found")
		}
		l.Logger.Error("[UpdatePaymentReceived] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query error: %v", err.Error())
	}

	// 校验是否属于当前用户
	if data.UserId != userId {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "payment received not found")
	}

	// 根据收款方式类型进行校验
	switch req.ReceivedType {
	case paymentreceived.ReceivedTypeBankcard:
		if req.ReceivedNo == "" {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "received_no is required for bankcard")
		}
		if req.BankName == "" {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "bank_name is required for bankcard")
		}
	case paymentreceived.ReceivedTypeWeixin, paymentreceived.ReceivedTypeAlipay:
		if req.ReceivedNo == "" {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "received_no is required")
		}
		if req.Qrcode == "" {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "qrcode is required for weixin/alipay")
		}
	default:
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "invalid received_type: %s", req.ReceivedType)
	}

	// 更新收款方式记录
	now := time.Now()
	data.ReceivedType = req.ReceivedType
	data.ReceivedNo = req.ReceivedNo
	data.BankName = req.BankName
	data.OpeningBranch = req.OpeningBranch
	data.Qrcode = req.Qrcode
	data.UpdatedAt = now

	err = l.svcCtx.PaymentReceivedModel.Update(l.ctx, data)
	if err != nil {
		l.Logger.Error("[UpdatePaymentReceived] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Update error: %v", err.Error())
	}

	return nil
}
