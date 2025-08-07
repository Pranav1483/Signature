package models

type ReqPayload struct {
	Head      Head      `json:"Head" xml:"Head"`
	Txn       Txn       `json:"Txn" xml:"Txn"`
	Signature Signature `json:"signature" xml:"Signature"`
}

type RespPayload struct {
	Head      Head      `json:"Head" xml:"Head"`
	Res       Res       `json:"res" xml:"Res"`
	Signature Signature `json:"signature" xml:"Signature"`
}

type ReqOnboard struct {
	Head      Head      `json:"Head" xml:"Head"`
	Onboard   Onboard   `json:"Onboard" xml:"Onboard"`
}

type RespOnboard struct {
	Head Head `json:"Head" xml:"Head"`
	Res  Res  `json:"Res" xml:"Res"`
}

type Head struct {
	RefNo     string `json:"refNo" xml:"refNo,attr"`
	Timestamp string `json:"timestamp" xml:"timestamp,attr"`
}

type Person struct {
	PersonType string `json:"personType" xml:"personType,attr"`
	Address    string `json:"address" xml:"address,attr"`
	OrgID      string `json:"orgId" xml:"orgId,attr"`
}

type Txn struct {
	TxnId    string `json:"txnId" xml:"txnId,attr"`
	Sender   Person `json:"sender" xml:"Sender"`
	Receiver Person `json:"receiver" xml:"Receiver"`
	Amount   string `json:"amount" xml:"Amount"`
}

type Res struct {
	Result  string `json:"result" xml:"result,attr"`
	ErrCode string `json:"errCode" xml:"errCode,attr"`
	Message string `json:"message" xml:"message,attr"`
	OrgId   string `json:"orgId" xml:"orgId,attr"`
}

type Onboard struct {
	OrgID           string `json:"orgId" xml:"orgId"`
	OrgName         string `json:"orgName" xml:"orgName"`
	URL             string `json:"url" xml:"url"`
	PublicKey       string `json:"publicKey" xml:"publicKey"`
	SignatureMethod string `json:"signatureMethod" xml:"signatureMethod"`
}

type Signature struct {
	Sign string `json:"sign" xml:"Sign"`
}

type Ack struct {
	RefNo     string `json:"refNo" xml:"refNo,attr"`
	Timestamp string `json:"timestamp" xml:"timestamp,attr"`
	Res       Res    `json:"res" xml:"res"`
}
