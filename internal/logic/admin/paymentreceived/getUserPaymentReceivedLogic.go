package paymentreceived

import (
	"context"

	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/logger"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
)

type GetUserPaymentReceivedLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPaymentReceivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPaymentReceivedLogic {
	return &GetUserPaymentReceivedLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPaymentReceivedLogic) GetUserPaymentReceived(req *types.GetUserPaymentReceivedRequest) (*types.GetUserPaymentReceivedResponse, error) {
	total, list, err := l.svcCtx.PaymentReceivedModel.QueryListByPage(l.ctx, int(req.Page), int(req.Size), req.UserId, req.ReceivedType)
	if err != nil {
		l.Logger.Error("[GetUserPaymentReceived] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Query error: %v", err.Error())
	}

	resp := &types.GetUserPaymentReceivedResponse{
		Total: total,
		List:  make([]types.PaymentReceived, 0, len(list)),
	}

	for _, item := range list {
		resp.List = append(resp.List, types.PaymentReceived{
			Id:            item.Id,
			UserId:        item.UserId,
			ReceivedName:  item.ReceivedName,
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
