package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"signature/internal/constants"

	ecies "github.com/ecies/go/v2"
)

type RSAService struct {
	ValidationType string
}

func (r *RSAService) Encrypt(data []byte, key string) ([]byte, error) {
	var result []byte
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, constants.ErrDecodePEMBlock
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	result, err = rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RSAService) Decrypt(data []byte, key string) ([]byte, error) {
	var result []byte
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, constants.ErrDecodePEMBlock
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	result, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *RSAService) ValidateSignature(data []byte, signature string, key string) (bool, error) {
	encrypted, err := r.Encrypt(data, key)
	if err != nil {
		return false, err
	}
	return string(encrypted) == signature, nil
}

func NewRSAService() *RSAService {
	return &RSAService{
		ValidationType: constants.RSA_SERVICE,
	}
}

type ECCService struct {
	ValidationType string
}

func (e *ECCService) Encrypt(data []byte, key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, constants.ErrDecodePEMBlock
	}
	pubIfc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("unable to parse PKIX public key")
	}
	pubKey := pubIfc.(*ecdsa.PublicKey)
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	eciesPublicKey, err := ecies.NewPublicKeyFromBytes(pubBytes)
	if err != nil {
		log.Println("85: ", err)
		return nil, err
	}
	encryptedData, err := ecies.Encrypt(eciesPublicKey, data)
	if err != nil {
		log.Println("91: ", err)
		return nil, err
	}
	return encryptedData, nil
}

func (e *ECCService) Decrypt(data []byte, key string) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, constants.ErrDecodePEMBlock
	}
	ecdsaKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privBytes := ecdsaKey.D.Bytes()
	eciesPrivateKey := ecies.NewPrivateKeyFromBytes(privBytes)
	return ecies.Decrypt(eciesPrivateKey, data)
}

func (e *ECCService) ValidateSignature(data []byte, signature string, key string) (bool, error) {
	encrypted, err := e.Encrypt(data, key)
	if err != nil {
		return false, err
	}
	return string(encrypted) == signature, nil
}

func NewECCService() *ECCService {
	return &ECCService{
		ValidationType: constants.ECC_SERVICE,
	}
}
