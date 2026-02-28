package paymentreceived

import (
	"context"

	"github.com/perfect-panel/server/internal/model/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/constant"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DeletePaymentReceivedLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePaymentReceivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePaymentReceivedLogic {
	return &DeletePaymentReceivedLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePaymentReceivedLogic) DeletePaymentReceived(req *types.DeletePaymentReceivedRequest) error {
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
		l.Logger.Error("[DeletePaymentReceived] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query error: %v", err.Error())
	}

	// 校验是否属于当前用户
	if data.UserId != userId {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "payment received not found")
	}

	// 执行硬删除
	err = l.svcCtx.PaymentReceivedModel.Delete(l.ctx, req.Id)
	if err != nil {
		l.Logger.Error("[DeletePaymentReceived] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Delete error: %v", err.Error())
	}

	return nil
}
