package paymentreceived

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/logic/admin/paymentreceived"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/result"
)

// Get User Payment Received
func GetUserPaymentReceivedHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.GetUserPaymentReceivedRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			result.ParamErrorResult(c, err)
			return
		}
		l := paymentreceived.NewGetUserPaymentReceivedLogic(c.Request.Context(), svcCtx)
		resp, err := l.GetUserPaymentReceived(&req)
		result.HttpResult(c, resp, err)
	}
}
