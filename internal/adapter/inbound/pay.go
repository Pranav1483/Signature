package inbound

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"signature/internal/constants"
	"signature/internal/entity/models"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Decrypt(c *gin.Context) ([]byte, string, error) {
	var (
		data    []byte
		err     error
		message string
	)
	reqCopy, err := io.ReadAll(c.Request.Body)
	if err != nil {
		message = "cannot read body"
		return data, message, err
	}
	signService, err := h.service.GetSignatureService(constants.SIGNMETHOD)
	if err != nil {
		message = "unable to get signature service"
		return data, message, err
	}
	data, err = signService.Decrypt(reqCopy, constants.PRIVATEKEY)
	if err != nil {
		message = "unable to decrypt"
		return data, message, err
	}
	return data, message, err
}

func (h *Handler) InitPay(c *gin.Context) {
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

func (h *Handler) ReqPay(c *gin.Context) {
	data, message, err := h.Decrypt(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}
	var reqPay models.ReqPayload
	if err := json.Unmarshal(data, &reqPay); err != nil {
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
	data, message, err := h.Decrypt(c)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}
	var respPay models.RespPayload
	if err := json.Unmarshal(data, &respPay); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid data"})
		return
	}

	go func() {
		h.worker <- func() {
			log.Println("implement")
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
