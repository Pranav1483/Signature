package models

type Organization struct {
	OrgID           string `json:"org_id" gorm:"column:org_id; primaryKey"`
	OrgName         string `json:"org_name" gorm:"column:org_name"`
	OrgType         string `json:"org_type" gorm:"column:org_type"`
	URL             string `json:"url" gorm:"column:url"`
	PublicKey       string `json:"public_key" gorm:"column:public_key"`
	SignatureMethod string `json:"signature_method" gorm:"column:signature_method"`
}
