// ============================================
// internal/repository/windows_certificate_repository.go
// OPTIMIZADO
// ============================================
package repository

import (
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"log"
	"server/internal/dto"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	crypt32                               = windows.NewLazySystemDLL("crypt32.dll")
	procCertGetCertificateContextProperty = crypt32.NewProc("CertGetCertificateContextProperty")
)

const (
	certKeyProvInfoPropID = 2
	storeName             = "MY"
)

type WindowsCertificateRepository struct{}

func NewWindowsCertificateRepository() *WindowsCertificateRepository {
	return &WindowsCertificateRepository{}
}

func (r *WindowsCertificateRepository) FindAll() ([]dto.Certificate, error) {
	log.Println("📂 Abriendo almacén de certificados...")

	store, err := windows.CertOpenSystemStore(0, windows.StringToUTF16Ptr(storeName))
	if err != nil {
		return nil, fmt.Errorf("error al abrir almacén de certificados: %w", err)
	}
	defer func() {
		log.Println("🔒 Cerrando almacén de certificados...")
		windows.CertCloseStore(store, 0)
	}()

	certificates := make([]dto.Certificate, 0, 32)
	seen := make(map[string]bool, 128)
	var ctx *windows.CertContext

	now := time.Now()

	log.Println("🔍 Enumerando certificados...")

	for {
		ctx, err = windows.CertEnumCertificatesInStore(store, ctx)
		if err != nil {
			if err == windows.ERROR_NO_MORE_ITEMS {
				log.Printf("✅ Enumeración completada. Total: %d certificados únicos\n", len(certificates))
				break
			}
			break
		}

		if ctx == nil {
			break
		}

		cert := r.parseCertContext(ctx, now)
		if cert == nil {
			continue
		}

		// Evitar duplicados por thumbprint
		if seen[cert.Thumbprint] {
			continue
		}
		seen[cert.Thumbprint] = true

		certificates = append(certificates, *cert)

		if len(certificates) >= 200 {
			log.Println("⚠️ Se alcanzó el límite de 200 certificados únicos (protección overflow)")
			break
		}
	}

	log.Printf("📊 Retornando %d certificados únicos\n", len(certificates))
	return certificates, nil
}

// ====================================================================================
// parseCertContext optimizado
// ====================================================================================
func (r *WindowsCertificateRepository) parseCertContext(ctx *windows.CertContext, now time.Time) *dto.Certificate {
	if ctx == nil || ctx.EncodedCert == nil || ctx.Length < 128 {
		return nil
	}

	// Leer el certificado de memoria nativa usando unsafe.Slice (mucho más rápido)
	certBytes := unsafe.Slice((*byte)(unsafe.Pointer(ctx.EncodedCert)), ctx.Length)

	parsed, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil
	}

	subject := parsed.Subject.String()
	issuer := parsed.Issuer.String()

	// Cache string lowercase para evitar ToLower en cada llamada
	subLower := strings.ToLower(subject)
	issLower := strings.ToLower(issuer)

	cert := &dto.Certificate{
		Subject:          subject,
		Issuer:           issuer,
		CommonName:       extractCommonName(subject),
		OwnerName:        extractOwnerName(subject),
		OrganizationName: extractOrganization(subject),
		IssuerName:       extractCommonName(issuer),
		DNI:              extractDNI(subject),
		SerialNumber:     parsed.SerialNumber.String(),
		NotBefore:        parsed.NotBefore,
		NotAfter:         parsed.NotAfter,
		HasPrivateKey:    r.hasPrivateKey(ctx),
		IsDNIe:           detectDNIe(subLower, issLower),
		Thumbprint:       fmt.Sprintf("%X", sha1.Sum(parsed.Raw)),
	}

	cert.IsValid = now.After(parsed.NotBefore) && now.Before(parsed.NotAfter)
	cert.DaysUntilExpiry = int(parsed.NotAfter.Sub(now).Hours() / 24)

	return cert
}

// ==============================================================================
// PRIVATE KEY CHECK
// ==============================================================================
func (r *WindowsCertificateRepository) hasPrivateKey(ctx *windows.CertContext) bool {
	if ctx == nil {
		return false
	}

	var size uint32

	ret, _, _ := procCertGetCertificateContextProperty.Call(
		uintptr(unsafe.Pointer(ctx)),
		certKeyProvInfoPropID,
		0,
		uintptr(unsafe.Pointer(&size)),
	)

	return ret != 0 && size > 0
}

// ==============================================================================
// EXTRACTION HELPERS (OPTIMIZADOS SIN REGEX)
// ==============================================================================
func extractCommonName(subject string) string {
	for _, part := range strings.Split(subject, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToUpper(part), "CN=") {
			return strings.TrimPrefix(part, "CN=")
		}
	}
	return ""
}

func extractOwnerName(subject string) string {
	cn := extractCommonName(subject)

	if idx := strings.Index(cn, "FAU"); idx > 0 {
		return strings.TrimSpace(cn[:idx])
	}

	return strings.TrimSuffix(cn, " soft")
}

func extractOrganization(subject string) string {
	for _, part := range strings.Split(subject, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToUpper(part), "O=") {
			return strings.TrimPrefix(part, "O=")
		}
	}
	return "CertSoft"
}

func extractDNI(subject string) string {
	// Mucho más rápido que regex
	// Busca: PNOPE-########
	idx := strings.Index(subject, "PNOPE-")
	if idx == -1 {
		return ""
	}

	if idx+6+8 <= len(subject) {
		return subject[idx+6 : idx+14]
	}
	return ""
}

// ==============================================================================
// DNIe DETECTION (OPTIMIZADO)
// ==============================================================================
func detectDNIe(subjectLower, issuerLower string) bool {
	// DNIe físico NO contiene "soft"
	if strings.Contains(subjectLower, "soft") {
		return false
	}

	hasPNOPE := strings.Contains(subjectLower, "pnope-")
	hasRENIEC := strings.Contains(issuerLower, "reniec")
	hasFAU := strings.Contains(subjectLower, "fau")

	return hasPNOPE && hasRENIEC && !hasFAU
}
