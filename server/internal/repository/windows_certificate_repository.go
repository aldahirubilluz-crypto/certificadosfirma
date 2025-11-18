// ============================================
// internal/repository/windows_certificate_repository.go
// ============================================
package repository

import (
	"crypto/sha1"
	"crypto/x509"
	"fmt"
	"log"
	"regexp"
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

	var certificates []dto.Certificate
	var ctx *windows.CertContext
	seen := make(map[string]bool)

	log.Println("🔍 Enumerando certificados...")

	for {
		ctx, err = windows.CertEnumCertificatesInStore(store, ctx)
		if err != nil {
			if err == windows.ERROR_NO_MORE_ITEMS {
				log.Printf("✅ Enumeración completada. Total: %d certificados únicos\n", len(certificates))
				break
			}
			log.Printf("⚠️ Error enumerando certificados: %v - Finalizando\n", err)
			break
		}

		if ctx == nil {
			log.Println("✅ Contexto nulo, finalizando enumeración")
			break
		}

		cert := r.parseCertContext(ctx)
		if cert != nil {
			if !seen[cert.Thumbprint] {
				seen[cert.Thumbprint] = true
				log.Printf("📜 #%d: %s (%s)\n", len(certificates)+1, cert.CommonName, cert.Tipo())
				certificates = append(certificates, *cert)
			}
		}

		if len(certificates) >= 100 {
			log.Println("⚠️ Límite de 100 certificados únicos alcanzado")
			break
		}
	}

	log.Printf("📊 Retornando %d certificados únicos\n", len(certificates))
	return certificates, nil
}

func (r *WindowsCertificateRepository) parseCertContext(ctx *windows.CertContext) *dto.Certificate {
	if ctx == nil || ctx.EncodedCert == nil || ctx.Length == 0 {
		return nil
	}

	certBytes := make([]byte, ctx.Length)
	for i := uint32(0); i < ctx.Length; i++ {
		certBytes[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(ctx.EncodedCert)) + uintptr(i)))
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil
	}

	thumbprint := fmt.Sprintf("%X", sha1.Sum(cert.Raw))
	commonName := extractCommonName(cert.Subject.String())
	ownerName := extractOwnerName(cert.Subject.String())
	orgName := extractOrganization(cert.Subject.String())
	issuerName := extractCommonName(cert.Issuer.String())
	dni := extractDNI(cert.Subject.String())
	hasPrivateKey := r.hasPrivateKey(ctx)
	isDNIe := detectDNIe(cert)

	now := time.Now()
	isValid := now.After(cert.NotBefore) && now.Before(cert.NotAfter)
	daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)

	return &dto.Certificate{
		Subject:          cert.Subject.String(),
		Issuer:           cert.Issuer.String(),
		CommonName:       commonName,
		OwnerName:        ownerName,
		OrganizationName: orgName,
		IssuerName:       issuerName,
		DNI:              dni,
		SerialNumber:     cert.SerialNumber.String(),
		NotBefore:        cert.NotBefore,
		NotAfter:         cert.NotAfter,
		HasPrivateKey:    hasPrivateKey,
		IsDNIe:           isDNIe,
		Thumbprint:       thumbprint,
		IsValid:          isValid,
		DaysUntilExpiry:  daysUntilExpiry,
	}
}

func (r *WindowsCertificateRepository) hasPrivateKey(ctx *windows.CertContext) bool {
	if ctx == nil {
		return false
	}

	var cbData uint32
	ret, _, _ := procCertGetCertificateContextProperty.Call(
		uintptr(unsafe.Pointer(ctx)),
		certKeyProvInfoPropID,
		0,
		uintptr(unsafe.Pointer(&cbData)),
	)
	return ret != 0 && cbData > 0
}

func extractCommonName(subject string) string {
	parts := strings.Split(subject, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToUpper(part), "CN=") {
			return strings.TrimPrefix(part, "CN=")
		}
	}
	return ""
}

func extractOwnerName(subject string) string {
	// Buscar el nombre completo en el CN
	cn := extractCommonName(subject)
	
	// Si es DNIe, extraer el nombre antes de "FAU" o el CN completo
	if strings.Contains(cn, "FAU") {
		parts := strings.Split(cn, "FAU")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	
	// Si tiene apellidos y nombres separados
	re := regexp.MustCompile(`CN=([^,]+)`)
	if matches := re.FindStringSubmatch(subject); len(matches) > 1 {
		name := strings.TrimSpace(matches[1])
		// Limpiar sufijos comunes
		name = strings.TrimSuffix(name, " soft")
		return name
	}
	
	return cn
}

func extractOrganization(subject string) string {
	parts := strings.Split(subject, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToUpper(part), "O=") {
			org := strings.TrimPrefix(part, "O=")
			return org
		}
	}
	return "CertSoft"
}

func extractDNI(subject string) string {
	re := regexp.MustCompile(`PNOPE-(\d{8})`)
	if matches := re.FindStringSubmatch(subject); len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func detectDNIe(cert *x509.Certificate) bool {
	subject := strings.ToLower(cert.Subject.String())
	issuer := strings.ToLower(cert.Issuer.String())

	// Un DNIe físico NUNCA tiene "soft" en el CN
	if strings.Contains(subject, "soft") {
		return false
	}

	// Un DNIe físico tiene PNOPE pero NO tiene "fau" (firma de autoridad)
	hasPNOPE := strings.Contains(subject, "pnope-")
	hasFAU := strings.Contains(strings.ToLower(cert.Subject.String()), "fau")
	hasRENIEC := strings.Contains(issuer, "reniec")

	// DNIe físico: tiene PNOPE, emitido por RENIEC, pero NO tiene FAU ni "soft"
	return hasPNOPE && hasRENIEC && !hasFAU
}