package kong

type Kong struct {
	Endpoint string
}

func New(endpoint string) *Kong {
	return &Kong{
		Endpoint: endpoint,
	}
}

func (k *Kong) UpdateCertificate() error {
	return nil
}
