package v1

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) getLogs(c *gin.Context) {
	logs, err := h.dashboardService.GetLogs(c.Request.Context())
	if err != nil {
		log.Println(err)
		writeErrorResponse(c, http.StatusInternalServerError, "logs", "failed to get logs")
		return
	}

	response := gin.H{"logs": logs}

	writeResponse(c, http.StatusOK, response)
}

func (h *Handler) deletePayment(c *gin.Context) {}

func (h *Handler) deleteLog(c *gin.Context) {}
