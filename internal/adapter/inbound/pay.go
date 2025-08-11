package inbound

import (
	"log"
	"net/http"
	"signature/internal/constants"
	"signature/internal/entity/models"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ReqPay(c *gin.Context) {
	var reqPay models.ReqPayload
	if err := c.BindJSON(&reqPay); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid data sent"})
		return
	}

	go func() {
		h.worker <- func() {
			h.service.GetPayService().HandleReqPay(reqPay, h.service.GetSignatureService)
		}
	}()

	ack := models.Ack{
		RefNo:     reqPay.Head.RefNo,
		Timestamp: time.Now().Format(time.RFC3339),
		Res: models.Res{
			Result:  constants.RESULT_SUCCESS,
			ErrCode: "",
			Message: "",
		},
	}
	c.JSON(http.StatusOK, ack)
}

func (h *Handler) RespPay(c *gin.Context) {
	var respPay models.RespPayload
	if err := c.BindJSON(&respPay); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid data"})
		return
	}

	go func() {
		h.worker <- func() {
			h.service.GetPayService().HandleRespPay(respPay, h.service.GetSignatureService)
		}
	}()

	ack := models.Ack{
		RefNo:     respPay.Head.RefNo,
		Timestamp: time.Now().Format(time.RFC3339),
		Res: models.Res{
			Result:  constants.RESULT_SUCCESS,
			ErrCode: "",
			Message: "",
		},
	}
	c.JSON(http.StatusOK, ack)
}
