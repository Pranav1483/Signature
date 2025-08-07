package service

import (
	"signature/internal/application/db"
	"signature/internal/entity/models"
)

type OrganizationService struct {
	db *db.DB
}

func (o *OrganizationService) GetOrganization(orgId string) (models.Organization, error) {
	return o.db.GetOrganization(orgId)
}

func (o *OrganizationService) GetOrganizationByURL(url string) (models.Organization, error) {
	return o.db.GetOrganizationByURL(url)
}

func (o *OrganizationService) SaveOrganization(orgId, orgName, orgType, url, publicKey, signatureMethod string) error {
	return o.db.SaveOrganization(orgId, orgName, orgType, url, publicKey, signatureMethod)
}

func NewOrganizationService(db *db.DB) *OrganizationService {
	return &OrganizationService{
		db: db,
	}
}
