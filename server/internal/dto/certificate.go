// ============================================
// internal/dto/certificate.go
// ============================================
package dto

import "time"

type Certificate struct {
	Subject          string
	Issuer           string
	CommonName       string
	OwnerName        string
	OrganizationName string
	IssuerName       string
	DNI              string
	SerialNumber     string
	NotBefore        time.Time
	NotAfter         time.Time
	HasPrivateKey    bool
	IsDNIe           bool
	Thumbprint       string
	IsValid          bool
	DaysUntilExpiry  int
}

// Tipo devuelve el tipo de certificado
func (c *Certificate) Tipo() string {
	if c.IsDNIe {
		return "DNIe"
	}
	if c.HasPrivateKey {
		return "Certificado"
	}
	return "CA"
}

// IsAutoFirmado verifica si el certificado es autofirmado
func (c *Certificate) IsAutoFirmado() bool {
	return c.Subject == c.Issuer
}

type CertificateOutput struct {
	ID                 string    `json:"id"`
	Nombre             string    `json:"nombre"`
	Tipo               string    `json:"tipo"`
	Issuer             string    `json:"issuer"`
	SerialNumber       string    `json:"serialNumber"`
	ValidoDesde        time.Time `json:"validoDesde"`
	ValidoHasta        time.Time `json:"validoHasta"`
	DiasParaVencer     int       `json:"diasParaVencer"`
	RequiereContrasena bool      `json:"requiereContrasena"`
	IsValid            bool      `json:"isValid"`
}

type CertificatesResponse struct {
	DniElectronico []CertificateOutput `json:"dniElectronico"`
	Certificados   []CertificateOutput `json:"certificados"`
}