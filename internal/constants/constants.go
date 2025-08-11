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
	RSA_SERVICE   = "RSA"
	ECDSA_SERVICE = "ECDSA"
)

const (
	PUBLIC_KEY      = "PUBLIC KEY"
	RSA_PRIVATE_KEY = "PRIVATE KEY"
	EC_PRIVATE_KEY  = "EC PRIVATE KEY"
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
