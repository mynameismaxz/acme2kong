package kong

import "github.com/mynameismaxz/acme2kong/pkg/logger"

type Kong struct {
	Endpoint string

	logger *logger.Logger
}

func New(endpoint string, log *logger.Logger) (*Kong, error) {
	return &Kong{
		Endpoint: endpoint,
		logger:   log,
	}, nil
}

func (k *Kong) UpdateCertificate(cert, privateKey []byte) error {
	return nil
}
