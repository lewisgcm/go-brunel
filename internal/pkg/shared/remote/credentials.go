package remote

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/pkg/errors"
)

// Credentials contain the Certificate and Key required for TLS mutual auth between servers and agents
type Credentials struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

const (
	commonName = "brunel"
)

// ServerConfig gets TLS configuration appropriate for a server
func (credentials *Credentials) ServerConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(credentials.Cert), []byte(credentials.Key))
	if err != nil {
		return nil, errors.Wrap(err, "error creating key pair from credentials")
	}

	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM([]byte(credentials.Cert)); !ok {
		return nil, errors.New("error appending certificate to pool")
	}

	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{cert},
		ClientAuth:               tls.RequireAndVerifyClientCert,
		ClientCAs:                clientCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}

// ClientConfig gets TLS configuration appropriate for an agent client
func (credentials *Credentials) ClientConfig() (*tls.Config, error) {
	clientCertPool := x509.NewCertPool()
	if ok := clientCertPool.AppendCertsFromPEM([]byte(credentials.Cert)); !ok {
		return nil, errors.New("error appending certificate to pool")
	}

	cert, err := tls.X509KeyPair([]byte(credentials.Cert), []byte(credentials.Key))
	if err != nil {
		return nil, errors.Wrap(err, "error creating key pair from credentials")
	}

	return &tls.Config{
		ServerName:   commonName,
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCertPool,
	}, nil
}

// GenerateCredentials will generate a new set of credentials for securing agent and server
func GenerateCredentials() (*Credentials, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2028)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate private key")
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate serial number")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 265),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	b, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create certificate")
	}

	return &Credentials{
		Cert: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: b})),
		Key:  string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})),
	}, nil
}
