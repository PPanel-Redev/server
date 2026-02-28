package paymentreceived

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/logic/admin/paymentreceived"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/pkg/result"
)

// List Payment Received By Type
func ListPaymentReceivedByTypeHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		l := paymentreceived.NewListPaymentReceivedByTypeLogic(c.Request.Context(), svcCtx)
		resp, err := l.ListPaymentReceivedByType()
		result.HttpResult(c, resp, err)
	}
}
