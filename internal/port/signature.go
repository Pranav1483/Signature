package port

type SignatureService interface {
	Generate(data []byte, privateKey string) ([]byte, error)
	Validate(data []byte, signature string, publicKey string) (bool, error)
}
