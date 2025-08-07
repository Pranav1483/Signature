package inbound

import (
	"net/http"
	"signature/internal/constants"
	"signature/internal/entity/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Onboard(c *gin.Context, orgType string) {
	var onboard models.ReqOnboard
	if err := c.BindJSON(&onboard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid data"})
		return
	}
	if err := h.service.GetOrganizationService().SaveOrganization(
		onboard.Onboard.OrgID, onboard.Onboard.OrgName, orgType, onboard.Onboard.URL,
		onboard.Onboard.PublicKey, onboard.Onboard.SignatureMethod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "err saving org"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "org onboarded"})
}

func (h *Handler) OnboardRegulator(c *gin.Context) {
	h.Onboard(c, constants.ORG_REGULATOR)
}

func (h *Handler) OnboardBank(c *gin.Context) {
	h.Onboard(c, constants.ORG_BANK)
}
