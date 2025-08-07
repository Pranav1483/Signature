package constants

import "errors"

var (
	ErrDecodePEMBlock              = errors.New("failed to decode PEM Block")
	ErrInvalidSignatureServiceType = errors.New("invalid signature service type")
)
