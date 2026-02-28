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
)

type ListPaymentReceivedLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPaymentReceivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPaymentReceivedLogic {
	return &ListPaymentReceivedLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPaymentReceivedLogic) ListPaymentReceived() (*types.GetPaymentReceivedListResponse, error) {
	u, ok := l.ctx.Value(constant.CtxKeyUser).(*user.User)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	userId := u.Id
	list, err := l.svcCtx.PaymentReceivedModel.FindListByUserId(l.ctx, userId)
	if err != nil {
		l.Logger.Error("[ListPaymentReceived] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query error: %v", err.Error())
	}

	resp := &types.GetPaymentReceivedListResponse{
		List: make([]types.PaymentReceived, 0, len(list)),
	}

	for _, item := range list {
		resp.List = append(resp.List, types.PaymentReceived{
			Id:            item.Id,
			UserId:        item.UserId,
			ReceivedType:  item.ReceivedType,
			ReceivedNo:    item.ReceivedNo,
			BankName:      item.BankName,
			OpeningBranch: item.OpeningBranch,
			Qrcode:        item.Qrcode,
			CreatedAt:     item.CreatedAt.Unix(),
		})
		if !item.UpdatedAt.IsZero() {
			resp.List[len(resp.List)-1].UpdatedAt = item.UpdatedAt.Unix()
		}
	}

	return resp, nil
}
