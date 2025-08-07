package constants

var (
	PORT       = ""
	PRIVATEKEY = ""
	SIGNMETHOD = ""
	ORGTYPE    = ""
	ORGID      = ""
)

var (
	DBHOST = ""
	DBPORT = ""
	DBNAME = ""
	DBUSER = ""
	DBPASS = ""
)

const (
	RSA_SERVICE = "RSA"
	ECC_SERVICE = "ECC"
)

const (
	RESULT_SUCCESS = "SUCCESS"
	RESULT_FAILURE = "FAILURE"
)

const (
	ORG_BANK      = "BANK"
	ORG_REGULATOR = "REGULATOR"
)

const (
	ENDPOINT_ONBOARDBANK      = "/Onboard/Bank"
	ENDPOINT_ONBOARDREGULATOR = "/Onboard/Regulator"
	ENDPOINT_REQPAY           = "/ReqPay"
	ENDPOINT_RESPPAY          = "/RespPay"
	ENDPOINT_PAY              = "/Pay"
)

const (
	CONTENT_JSON = "application/json"
)
