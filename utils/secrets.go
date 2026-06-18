package utils

import (
	"context"
	"crypto/rsa"
	"log/slog"
	"os"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

//var SecretKey = []byte("secretkey")

type Secret struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	FilePath   string
	Password   string
	logger     *slog.Logger
}

func NewSecret(path, password string, logger *slog.Logger) *Secret {
	return &Secret{
		FilePath: path,
		Password: password,
	}
}

func (s *Secret) Extract() (*rsa.PrivateKey, *rsa.PublicKey) {
	fileContent, err := os.ReadFile(s.FilePath)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelError,
			"exception in loading the secret",
			slog.String("error", err.Error()))
	}

	privateKey, certificate, err := pkcs12.Decode(fileContent, s.Password)
	if err != nil {
		s.logger.LogAttrs(context.Background(), slog.LevelError,
			"exception in decoding the privatekey and certificate",
			slog.String("error", err.Error()))
	}

	s.PrivateKey = privateKey.(*rsa.PrivateKey)
	s.PublicKey = certificate.PublicKey.(*rsa.PublicKey)
	return s.PrivateKey, s.PublicKey
}
