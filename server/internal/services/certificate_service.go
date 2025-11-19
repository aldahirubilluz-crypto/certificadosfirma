// ============================================
// internal/services/certificate_service.go
// ============================================
package services

import (
	"server/internal/dto"
	"server/internal/repository"
)

type CertificateService struct {
	repo repository.CertificateRepository
}

func NewCertificateService() *CertificateService {
	return &CertificateService{
		repo: repository.NewWindowsCertificateRepository(),
	}
}

func (s *CertificateService) GetCertificatesGrouped() (*dto.CertificatesResponse, error) {
	allCerts, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	dniElectronicos := []dto.CertificateOutput{}
	certificados := []dto.CertificateOutput{}

	for _, cert := range allCerts {
		// Solo procesar certificados con clave privada
		if !cert.HasPrivateKey {
			continue
		}

		// Filtrar certificados autofirmados (issuer == subject)
		if cert.IsAutoFirmado() {
			continue
		}

		output := dto.CertificateOutput{
			ID:                 cert.Thumbprint,
			Nombre:             cert.CommonName,
			Issuer:             cert.IssuerName,
			SerialNumber:       cert.SerialNumber,
			ValidoDesde:        cert.NotBefore,
			ValidoHasta:        cert.NotAfter,
			DiasParaVencer:     cert.DaysUntilExpiry,
			RequiereContrasena: cert.HasPrivateKey,
			IsValid:            cert.IsValid,
		}

		if cert.IsDNIe {
			output.Tipo = "DNIe"
			output.Nombre = cert.OwnerName + " - DNIe"
			dniElectronicos = append(dniElectronicos, output)
		} else {
			output.Tipo = "Certificado"
			output.Nombre = cert.OwnerName + " - " + cert.OrganizationName
			certificados = append(certificados, output)
		}
	}

	return &dto.CertificatesResponse{
		DniElectronico: dniElectronicos,
		Certificados:   certificados,
	}, nil
}
