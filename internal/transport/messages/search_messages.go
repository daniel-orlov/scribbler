package messages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) SearchMessages(gctx *gin.Context) {
	h.logger.Debug("endpoint hit", zap.String("handler", "SearchMessages"))

	gctx.JSON(http.StatusCreated, nil)
}
