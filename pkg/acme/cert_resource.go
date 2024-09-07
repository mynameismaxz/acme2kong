package acme

import (
	"encoding/json"

	"github.com/go-acme/lego/v4/certificate"
)

type CertResource struct {
	Domain            string `json:"domain"`
	CertURL           string `json:"certUrl"`
	CertStableURL     string `json:"certStableUrl"`
	PrivateKey        []byte `json:"privateKey"`
	Certificate       []byte `json:"certificate"`
	IssuerCertificate []byte `json:"issuerCertificate"`
	CSR               []byte `json:"csr"`
}

func LoadCertResourceFromBytes(data []byte) (*CertResource, error) {
	var result CertResource
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func ConvertToCertResource(cr *certificate.Resource) *CertResource {
	var result CertResource
	result.Domain = cr.Domain
	result.CertURL = cr.CertURL
	result.CertStableURL = cr.CertStableURL
	result.PrivateKey = cr.PrivateKey
	result.Certificate = cr.Certificate
	result.IssuerCertificate = cr.IssuerCertificate
	result.CSR = cr.CSR
	return &result
}

func (cr *CertResource) toBytes() []byte {
	bytes, err := json.Marshal(cr)
	if err != nil {
		// handle error
	}
	return bytes
}

func (cr *CertResource) toACMECertificateResource() *certificate.Resource {
	var result certificate.Resource
	result.Domain = cr.Domain
	result.CertURL = cr.CertURL
	result.CertStableURL = cr.CertStableURL
	result.PrivateKey = cr.PrivateKey
	result.Certificate = cr.Certificate
	result.IssuerCertificate = cr.IssuerCertificate
	result.CSR = cr.CSR
	return &result
}
