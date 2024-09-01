package kong

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"path"

	"github.com/mynameismaxz/acme2kong/pkg/httpclient"
	"github.com/mynameismaxz/acme2kong/pkg/logger"
)

const (
	certificates = "certificates"
)

type Kong struct {
	Endpoint   string
	DomainName []string

	logger *logger.Logger
}

type Certificate struct {
	Cert string   `json:"cert"`
	Key  string   `json:"key"`
	Snis []string `json:"snis"`
}

func New(endpoint string, domain []string, log *logger.Logger) (*Kong, error) {
	return &Kong{
		Endpoint:   endpoint,
		DomainName: domain,
		logger:     log,
	}, nil
}

func (k *Kong) UpdateCertificate(cert, privateKey []byte) error {

	endPoint := path.Clean(path.Join(k.Endpoint, certificates))
	k.DomainName = []string{"test.mymacz.com"}

	certJSON, err := json.Marshal(Certificate{
		Cert: string(cert),
		Key:  string(privateKey),
		Snis: k.DomainName,
	})
	if err != nil {
		return err
	}

	// http request to kong post
	client, err := http.NewRequest(http.MethodPost, endPoint, bytes.NewBuffer(certJSON))
	if err != nil {
		return err
	}
	client.Header.Set("Content-Type", "application/json")

	resp, err := httpclient.Dorequest(client)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		k.logger.Error("Bad request")
		return errors.New("bad request")
	}

	return nil
}
