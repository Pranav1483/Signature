package port

type SignatureService interface {
	Encrypt(data []byte, publicKey string) ([]byte, error)
	Decrypt(data []byte, privateKey string) ([]byte, error)
	ValidateSignature(data []byte, signature string, publicKey string) (bool, error)
}
