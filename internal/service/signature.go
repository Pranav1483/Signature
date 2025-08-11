package service

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"log"
	"math/big"
	"signature/internal/constants"
)

type RSAService struct {
	ValidationType string
}

func (r *RSAService) Generate(data []byte, key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != constants.RSA_PRIVATE_KEY {
		return nil, constants.ErrDecodePEMBlock
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(data)
	signature, err := rsa.SignPSS(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hash[:], nil)
	if err != nil {
		return nil, err
	}

	sigB64 := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(sigB64, signature)

	return sigB64, nil
}

func (r *RSAService) Validate(data []byte, signature string, key string) (bool, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != constants.PUBLIC_KEY {
		return false, constants.ErrDecodePEMBlock
	}
	pubIface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	pub, ok := pubIface.(*rsa.PublicKey)
	if !ok {
		return false, constants.ErrRSAKey
	}

	sig, err := base64.StdEncoding.DecodeString(string(signature))
	if err != nil {
		return false, err
	}

	hash := sha256.Sum256(data)

	err = rsa.VerifyPSS(pub, crypto.SHA256, hash[:], sig, nil)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func NewRSAService() *RSAService {
	return &RSAService{
		ValidationType: constants.RSA_SERVICE,
	}
}

type ECDSAService struct {
	ValidationType string
}

type asn1Sig struct {
	R, S *big.Int
}

func (e *ECDSAService) Generate(data []byte, key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != constants.EC_PRIVATE_KEY {
		return nil, constants.ErrDecodePEMBlock
	}
	priv, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return nil, err
	}
	sigBytes, err := asn1.Marshal(asn1Sig{R: r, S: s})
	if err != nil {
		return nil, err
	}
	sigB4 := make([]byte, base64.StdEncoding.EncodedLen(len(sigBytes)))
	base64.StdEncoding.Encode(sigB4, sigBytes)
	return sigB4, nil
}

func (e *ECDSAService) Validate(data []byte, signature string, key string) (bool, error) {
	block, _ := pem.Decode([]byte(key))
	log.Println(block)
	if block == nil || block.Type != constants.PUBLIC_KEY {
		return false, constants.ErrDecodePEMBlock
	}
	pubIface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	pub, ok := pubIface.(*ecdsa.PublicKey)
	if !ok {
		return false, constants.ErrECDSAKey
	}
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	var sig asn1Sig
	if _, err = asn1.Unmarshal(sigBytes, &sig); err != nil {
		return false, err
	}
	hash := sha256.Sum256(data)
	return ecdsa.Verify(pub, hash[:], sig.R, sig.S), nil
}

func NewECDSAService() *ECDSAService {
	return &ECDSAService{
		ValidationType: constants.ECDSA_SERVICE,
	}
}

type EDDSAService struct {
	ValidationType string
}

func (e *EDDSAService) Generate(data []byte, key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != constants.ED_PRIVATE_KEY {
		return nil, constants.ErrDecodePEMBlock
	}
	keyIface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := keyIface.(ed25519.PrivateKey)
	if !ok {
		return nil, constants.ErrED25519Key
	}
	sig := ed25519.Sign(priv, data)
	sigB64 := make([]byte, base64.StdEncoding.EncodedLen(len(sig)))
	base64.StdEncoding.Encode(sigB64, sig)
	return sigB64, nil
}

func (e *EDDSAService) Validate(data []byte, signature string, key string) (bool, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil || block.Type != constants.PUBLIC_KEY {
		return false, constants.ErrDecodePEMBlock
	}
	pubIface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	pub, ok := pubIface.(ed25519.PublicKey)
	if !ok {
		return false, constants.ErrED25519Key
	}
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	valid := ed25519.Verify(pub, data, sig)
	return valid, nil
}

func NewEDDSAService() *EDDSAService {
	return &EDDSAService{
		ValidationType: constants.EDDSA_SERVICE,
	}
}
