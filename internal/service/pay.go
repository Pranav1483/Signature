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
		if err = p.HandleReqPayBank(reqPay, GetSignatureService); err != nil {
			log.Println("BANK:", err)
		}
	case constants.ORG_REGULATOR:
		if err = p.HandleReqPayRegulator(reqPay, GetSignatureService); err != nil {
			log.Println("REGULATOR:", err)
		}
	default:
		return
	}
}

func (p *PayService) HandleReqPayBank(reqPay models.ReqPayload, GetSignatureService func(string) (port.SignatureService, error)) error {
	var (
		data     []byte
		err      error
		endpoint string
	)
	regulator, err := p.db.GetOrganizationByType(constants.ORG_REGULATOR)
	if err != nil {
		return err
	}
	if reqPay.Txn.Sender.OrgID == constants.ORGID {
		endpoint = constants.ENDPOINT_REQPAY
		data, err = json.Marshal(reqPay)
		if err != nil {
			return err
		}
		signatureService, err := GetSignatureService(constants.SIGNMETHOD)
		if err != nil {
			return err
		}
		signature, err := signatureService.Generate(data, constants.PRIVATEKEY)
		if err != nil {
			return err
		}
		reqPay.Signature = models.Signature{
			Sign: string(signature),
		}
		data, err = json.Marshal(reqPay)
		if err != nil {
			return err
		}
	} else {
		endpoint = constants.ENDPOINT_RESPPAY
		signature := models.Signature{
			Sign: reqPay.Signature.Sign,
		}
		reqPay.Signature = models.Signature{}
		data, err = json.Marshal(reqPay)
		if err != nil {
			return err
		}
		signatureService, err := GetSignatureService(regulator.SignatureMethod)
		if err != nil {
			return err
		}
		valid, err := signatureService.Validate(data, signature.Sign, regulator.PublicKey)
		if !valid {
			return constants.ErrInvalidSignature
		} else if err != nil {
			return err
		}
		respPay := models.RespPayload{
			Head: models.Head{
				RefNo:     reqPay.Head.RefNo,
				Timestamp: time.Now().Format(time.RFC3339),
			},
			Res: models.Res{
				Result:        constants.RESULT_SUCCESS,
				ErrCode:       "",
				Message:       "",
				SenderOrgId:   reqPay.Txn.Receiver.OrgID,
				ReceiverOrgId: reqPay.Txn.Sender.OrgID,
			},
		}
		data, err = json.Marshal(respPay)
		if err != nil {
			return err
		}
		signatureService, err = GetSignatureService(constants.SIGNMETHOD)
		if err != nil {
			return err
		}
		newSignature, err := signatureService.Generate(data, constants.PRIVATEKEY)
		if err != nil {
			return err
		}
		respPay.Signature = models.Signature{
			Sign: string(newSignature),
		}
		data, err = json.Marshal(respPay)
		if err != nil {
			return err
		}
	}
	http.Post(fmt.Sprintf("%s%s", regulator.URL, endpoint), constants.CONTENT_JSON, bytes.NewReader(data))
	return nil
}

func (p *PayService) HandleReqPayRegulator(reqPay models.ReqPayload, GetSignatureService func(string) (port.SignatureService, error)) error {
	senderOrg, err := p.db.GetOrganization(reqPay.Txn.Sender.OrgID)
	if err != nil {
		return err
	}
	receiverOrg, err := p.db.GetOrganization(reqPay.Txn.Receiver.OrgID)
	if err != nil {
		return err
	}
	signature := models.Signature{
		Sign: reqPay.Signature.Sign,
	}
	reqPay.Signature = models.Signature{}
	data, err := json.Marshal(reqPay)
	if err != nil {
		return err
	}
	signatureService, err := GetSignatureService(senderOrg.SignatureMethod)
	if err != nil {
		return err
	}
	valid, err := signatureService.Validate(data, signature.Sign, senderOrg.PublicKey)
	if !valid {
		log.Println(err)
		return constants.ErrInvalidSignature
	} else if err != nil {
		return err
	}
	signatureService, err = GetSignatureService(constants.SIGNMETHOD)
	if err != nil {
		return err
	}
	newSignature, err := signatureService.Generate(data, constants.PRIVATEKEY)
	if err != nil {
		return err
	}
	reqPay.Signature = models.Signature{
		Sign: string(newSignature),
	}
	data, err = json.Marshal(reqPay)
	if err != nil {
		return err
	}
	http.Post(fmt.Sprintf("%s%s", receiverOrg.URL, constants.ENDPOINT_REQPAY), constants.CONTENT_JSON, bytes.NewReader(data))
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
	senderOrg, err := p.db.GetOrganization(respPay.Res.SenderOrgId)
	if err != nil {
		return err
	}
	receiverOrg, err := p.db.GetOrganization(respPay.Res.ReceiverOrgId)
	if err != nil {
		return err
	}
	signature := models.Signature{
		Sign: respPay.Signature.Sign,
	}
	respPay.Signature = models.Signature{}
	data, err := json.Marshal(respPay)
	if err != nil {
		return err
	}
	signatureService, err := GetSignatureService(senderOrg.SignatureMethod)
	if err != nil {
		return err
	}
	valid, err := signatureService.Validate(data, signature.Sign, senderOrg.PublicKey)
	if !valid {
		return constants.ErrInvalidSignature
	} else if err != nil {
		return err
	}
	signatureService, err = GetSignatureService(constants.SIGNMETHOD)
	if err != nil {
		return err
	}
	newSignature, err := signatureService.Generate(data, constants.PRIVATEKEY)
	if err != nil {
		return err
	}
	respPay.Signature.Sign = string(newSignature)
	data, err = json.Marshal(respPay)
	if err != nil {
		return err
	}
	http.Post(fmt.Sprintf("%s%s", receiverOrg.URL, constants.ENDPOINT_RESPPAY), constants.CONTENT_JSON, bytes.NewReader(data))
	return nil
}

func NewPayService(db *db.DB) *PayService {
	return &PayService{
		db: db,
	}
}
