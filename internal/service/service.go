package service

import (
	"signature/internal/application/db"
	"signature/internal/constants"
	"signature/internal/port"
)

type Service struct {
	organizationService *OrganizationService
	payService          *PayService
	eccService          *ECCService
	rsaService          *RSAService
}

func (s *Service) GetSignatureService(serviceType string) (port.SignatureService, error) {
	switch serviceType {
	case constants.RSA_SERVICE:
		return s.rsaService, nil
	case constants.ECC_SERVICE:
		return s.eccService, nil
	default:
		return nil, constants.ErrInvalidSignatureServiceType
	}
}

func (s *Service) GetOrganizationService() *OrganizationService {
	return s.organizationService
}

func (s *Service) GetPayService() *PayService {
	return s.payService
}

func NewService(db *db.DB) *Service {
	return &Service{
		eccService:          NewECCService(),
		rsaService:          NewRSAService(),
		organizationService: NewOrganizationService(db),
		payService:          NewPayService(db),
	}
}
