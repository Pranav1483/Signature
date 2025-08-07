package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"signature/internal/application/db"
	"signature/internal/constants"
	"signature/internal/entity/models"
	"signature/internal/port"
	"time"
)

type PayService struct {
	db *db.DB
}

func (p *PayService) HandleReqPay(reqPay models.ReqPayload, GetSignatureService func(string) (port.SignatureService, error)) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	switch constants.ORGTYPE {
	case constants.ORG_BANK:
		err = p.HandleReqPayBank(reqPay, GetSignatureService)
		log.Println("BANK: ", err)
	case constants.ORG_REGULATOR:
		err = p.HandleReqPayRegulator(reqPay, GetSignatureService)
		log.Println("REGULATOR: ", err)
	default:
		return
	}
}

func (p *PayService) HandleReqPayBank(reqPay models.ReqPayload, GetSignatureService func(string) (port.SignatureService, error)) error {
	var (
		data []byte
		err  error
	)
	if reqPay.Txn.Sender.OrgID == constants.ORGID {
		data, err = json.Marshal(reqPay)
		if err != nil {
			return err
		}
	} else {
		respPay := models.RespPayload{
			Head: models.Head{
				RefNo:     reqPay.Head.RefNo,
				Timestamp: time.Now().Format(time.RFC3339),
			},
			Res: models.Res{
				Result:  constants.RESULT_SUCCESS,
				ErrCode: "",
				Message: "",
				OrgId:   reqPay.Txn.Receiver.OrgID,
			},
			Signature: models.Signature{Sign: ""},
		}
		data, err = json.Marshal(respPay)
		if err != nil {
			return err
		}
	}

	regulator, err := p.db.GetOrganizationByType(constants.ORG_REGULATOR)
	if err != nil {
		return err
	}
	signatureService, err := GetSignatureService(regulator.SignatureMethod)
	if err != nil {
		return err
	}
	encryptedData, err := signatureService.Encrypt(data, regulator.PublicKey)
	if err != nil {
		return err
	}
	http.Post(fmt.Sprintf("%s%s", regulator.URL, constants.ENDPOINT_RESPPAY), constants.CONTENT_JSON, bytes.NewReader(encryptedData))
	return nil
}

func (p *PayService) HandleReqPayRegulator(reqPay models.ReqPayload, GetSignatureService func(string) (port.SignatureService, error)) error {
	receiverOrg, err := p.db.GetOrganization(reqPay.Txn.Receiver.OrgID)
	if err != nil {
		return err
	}
	data, err := json.Marshal(reqPay)
	if err != nil {
		return err
	}
	signatureService, err := GetSignatureService(receiverOrg.SignatureMethod)
	if err != nil {
		return err
	}
	encryptedData, err := signatureService.Encrypt(data, receiverOrg.PublicKey)
	if err != nil {
		return err
	}
	http.Post(fmt.Sprintf("%s%s", receiverOrg.URL, constants.ENDPOINT_RESPPAY), constants.CONTENT_JSON, bytes.NewReader(encryptedData))
	return nil
}

func (p *PayService) HandleRespPay(respPay models.RespPayload, GetSignatureService func(string) (port.SignatureService, error)) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	switch constants.ORGTYPE {
	case constants.ORG_BANK:
		err = p.HandleRespPayBank(respPay, GetSignatureService)
	case constants.ORG_REGULATOR:
		err = p.HandleRespPayRegulator(respPay, GetSignatureService)
	default:
		return
	}
}

func (p *PayService) HandleRespPayBank(respPay models.RespPayload, GetSignatureService func(string) (port.SignatureService, error)) error {
	log.Printf("Txn %s done\n", respPay.Head.RefNo)
	return nil
}

func (p *PayService) HandleRespPayRegulator(respPay models.RespPayload, GetSignatureService func(string) (port.SignatureService, error)) error {
	receiverOrg, err := p.db.GetOrganization(respPay.Res.OrgId)
	if err != nil {
		return err
	}
	data, err := json.Marshal(respPay)
	if err != nil {
		return err
	}
	signatureService, err := GetSignatureService(receiverOrg.SignatureMethod)
	if err != nil {
		return err
	}
	encryptedData, err := signatureService.Encrypt(data, receiverOrg.PublicKey)
	if err != nil {
		return err
	}
	http.Post(fmt.Sprintf("%s%s", receiverOrg.URL, constants.ENDPOINT_RESPPAY), constants.CONTENT_JSON, bytes.NewReader(encryptedData))
	return nil
}

func NewPayService(db *db.DB) *PayService {
	return &PayService{
		db: db,
	}
}
