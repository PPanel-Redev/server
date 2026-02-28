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
)

type CreatePaymentReceivedLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePaymentReceivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentReceivedLogic {
	return &CreatePaymentReceivedLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaymentReceivedLogic) CreatePaymentReceived(req *types.CreatePaymentReceivedRequest) error {
	u, ok := l.ctx.Value(constant.CtxKeyUser).(*user.User)
	if !ok {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	userId := u.Id
	// 校验收款人姓名
	if req.ReceivedName == "" {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "received_name is required")
	}
	// 根据收款方式类型进行校验
	switch req.ReceivedType {
	case paymentreceived.ReceivedTypeBankcard:
		// 银行卡：received_no 和 bank_name 必填，qrcode 选填
		if req.ReceivedNo == "" {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "received_no is required for bankcard")
		}
		if req.BankName == "" {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "bank_name is required for bankcard")
		}
	case paymentreceived.ReceivedTypeWeixin, paymentreceived.ReceivedTypeAlipay:
		// 微信/支付宝：received_no 必填，qrcode 选填，bank_name 和 opening_branch 不需要填写
		if req.ReceivedNo == "" {
			return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "received_no is required")
		}
	default:
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "invalid received_type: %s", req.ReceivedType)
	}

	// 创建收款方式记录
	data := &paymentreceived.PaymentReceived{
		UserId:        userId,
		ReceivedName:  req.ReceivedName,
		ReceivedType:  req.ReceivedType,
		ReceivedNo:    req.ReceivedNo,
		BankName:      req.BankName,
		OpeningBranch: req.OpeningBranch,
		Qrcode:        req.Qrcode,
		CreatedAt:     time.Now(),
		IsDel:         1,
	}

	err := l.svcCtx.PaymentReceivedModel.Insert(l.ctx, data)
	if err != nil {
		l.Logger.Error("[CreatePaymentReceived] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "Insert error: %v", err.Error())
	}

	return nil
}
