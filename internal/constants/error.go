package constants

import "errors"

var (
	ErrDecodePEMBlock              = errors.New("failed to decode PEM Block")
	ErrInvalidSignatureServiceType = errors.New("invalid signature service type")
	ErrRSAKey                      = errors.New("not an RSA Key")
	ErrECDSAKey                      = errors.New("not an ECDSA Key")
	ErrInvalidSignature            = errors.New("invalid signature")
)
