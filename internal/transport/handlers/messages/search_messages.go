package messages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) SearchMessages(gctx *gin.Context) {
	h.logger.Debug("endpoint hit", zap.String("handler", "SearchMessages"))

	err := h.uc.FetchMessages(gctx.Request.Context(), nil)
	if err != nil {
		h.logger.Error("failed to fetch messages", zap.Error(err))
		gctx.JSON(http.StatusInternalServerError, nil)

		return
	}

	gctx.JSON(http.StatusCreated, nil)
}
