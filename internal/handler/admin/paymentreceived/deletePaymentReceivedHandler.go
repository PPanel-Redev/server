package paymentreceived

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/logic/admin/paymentreceived"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/result"
)

// Delete Payment Received
func DeletePaymentReceivedHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.DeletePaymentReceivedRequest
		_ = c.ShouldBind(&req)
		validateErr := svcCtx.Validate(&req)
		if validateErr != nil {
			result.ParamErrorResult(c, validateErr)
			return
		}

		l := paymentreceived.NewDeletePaymentReceivedLogic(c.Request.Context(), svcCtx)
		err := l.DeletePaymentReceived(&req)
		result.HttpResult(c, nil, err)
	}
}
