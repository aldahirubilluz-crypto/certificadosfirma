// ============================================
// internal/repository/certificate_repository.go
// ============================================
package repository

import "server/internal/dto"

type CertificateRepository interface {
	FindAll() ([]dto.Certificate, error)
}