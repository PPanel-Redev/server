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

type ListPaymentReceivedByTypeLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPaymentReceivedByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPaymentReceivedByTypeLogic {
	return &ListPaymentReceivedByTypeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPaymentReceivedByTypeLogic) ListPaymentReceivedByType() (*types.GetPaymentReceivedByTypeResponse, error) {
	u, ok := l.ctx.Value(constant.CtxKeyUser).(*user.User)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	userId := u.Id
	list, err := l.svcCtx.PaymentReceivedModel.FindListByUserId(l.ctx, userId)
	if err != nil {
		l.Logger.Error("[ListPaymentReceivedByType] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query error: %v", err.Error())
	}

	// Group by received_type
	typeMap := make(map[string][]types.PaymentReceived)
	for _, item := range list {
		payment := types.PaymentReceived{
			Id:            item.Id,
			UserId:        item.UserId,
			ReceivedType:  item.ReceivedType,
			ReceivedNo:    item.ReceivedNo,
			BankName:      item.BankName,
			OpeningBranch: item.OpeningBranch,
			Qrcode:        item.Qrcode,
			CreatedAt:     item.CreatedAt.Unix(),
		}
		if !item.UpdatedAt.IsZero() {
			payment.UpdatedAt = item.UpdatedAt.Unix()
		}
		typeMap[item.ReceivedType] = append(typeMap[item.ReceivedType], payment)
	}

	resp := &types.GetPaymentReceivedByTypeResponse{
		ReceivedType: make([]types.PaymentReceivedByType, 0, len(typeMap)),
	}

	for receivedType, items := range typeMap {
		resp.ReceivedType = append(resp.ReceivedType, types.PaymentReceivedByType{
			Type: receivedType,
			List: items,
		})
	}

	return resp, nil
}
