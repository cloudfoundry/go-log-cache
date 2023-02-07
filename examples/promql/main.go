package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	client "code.cloudfoundry.org/go-log-cache/v2"

	envstruct "code.cloudfoundry.org/go-envstruct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <query>", os.Args[0])
	}

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("invalid configuration: %s", err)
	}

	client := client.NewClient(
		cfg.LogCacheAddr,
		client.WithViaGRPC(
			grpc.WithTransportCredentials(cfg.TLS.Credentials("log-cache")),
		),
	)

	result, err := client.PromQL(context.Background(), os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	str, err := protojson.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(str))
}

// Config is the configuration for a LogCache Gateway.
type Config struct {
	LogCacheAddr string `env:"LOG_CACHE_ADDR, required"`
	TLS          TLS
}

type TLS struct {
	CAPath   string `env:"CA_PATH,   required"`
	CertPath string `env:"CERT_PATH, required"`
	KeyPath  string `env:"KEY_PATH,  required"`
}

func (t TLS) Credentials(cn string) credentials.TransportCredentials {
	creds, err := NewTLSCredentials(t.CAPath, t.CertPath, t.KeyPath, cn)
	if err != nil {
		log.Fatalf("failed to load TLS config: %s", err)
	}

	return creds
}

func NewTLSCredentials(
	caPath string,
	certPath string,
	keyPath string,
	cn string,
) (credentials.TransportCredentials, error) {
	cfg, err := NewTLSConfig(caPath, certPath, keyPath, cn)
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(cfg), nil
}

func NewTLSConfig(caPath, certPath, keyPath, cn string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ServerName:         cn,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
	}

	caCertBytes, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertBytes); !ok {
		return nil, errors.New("cannot parse ca cert")
	}

	tlsConfig.RootCAs = caCertPool

	return tlsConfig, nil
}

func LoadConfig() (*Config, error) {
	c := Config{}

	if err := envstruct.Load(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
