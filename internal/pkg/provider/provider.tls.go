// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"sync/atomic"

	"github.com/searKing/sole/pkg/database/models/key"
	"github.com/searKing/sole/pkg/database/orm"
	"gopkg.in/square/go-jose.v2"

	"github.com/pkg/errors"
	tls_ "github.com/searKing/golang/go/crypto/tls"
	"golang.org/x/net/http2"
)

func tlsKeyName() string {
	return GlobalProvider().Proto().GetService().GetName() + ".https-tls"
}

//go:generate go-atomicvalue -type "tlsConfig<*crypto/tls.Config>"
type tlsConfig atomic.Value

func (p *Provider) TLSConfig() *tls.Config {
	return p.tlsConfig.Load()
}

func (p *Provider) updateTLSConfig() {
	proto := p.Proto()
	logger := GlobalProvider().Logger().WithField("module", "provider.tls")

	var tlsConfig *tls.Config
	if proto.GetWeb().GetForceDisableTls() {
		logger.Warnln("HTTPS disabled. Never do this in production.")
	} else if len(proto.GetWeb().GetTls().GetAllowedTlsCidrs()) != 0 {
		logger.Infoln("TLS termination enabled, disabling https.")
	} else {
		tlsConfig = GetOrCreateTLSConfig()
		tlsConfig.InsecureSkipVerify = true
		tlsConfig.NextProtos = append(tlsConfig.NextProtos, http2.NextProtoTLS)
	}
	p.tlsConfig.Store(tlsConfig)
}

func GetOrCreateTLSConfig() *tls.Config {
	certs, cas := GetOrCreateTLSCertificate()
	return &tls.Config{
		ServerName:   GlobalProvider().Proto().GetWeb().GetTls().GetServiceName(),
		RootCAs:      cas,
		Certificates: certs,
	}

}

func GetOrCreateTLSCertificate() ([]tls.Certificate, *x509.CertPool) {
	proto := GlobalProvider().Proto()
	logger := GlobalProvider().Logger().WithField("module", "provider.tls")
	keyManager := GlobalProvider().KeyManager()
	tlsProto := proto.GetWeb().GetTls()

	certs, certPool, err := tls_.LoadCertificateAndPool(nil,
		tlsProto.GetKeyPairBase64().GetCert(), tlsProto.GetKeyPairBase64().GetKey(),
		tlsProto.GetKeyPairPath().GetCert(), tlsProto.GetKeyPairPath().GetKey())
	if err == nil {
		return certs, certPool
	}
	if errors.Cause(err) != tls_.ErrNoCertificatesConfigured {
		logger.WithError(errors.WithStack(err)).Fatalf("Unable to load HTTPS TLS Certificate")
	}

	_, privateKey, err := AsymmetricKeypair(GlobalProvider().Context(), &key.RS256Generator{KeyLength: 4069}, tlsKeyName())
	if err != nil {
		logger.WithError(errors.WithStack(err)).Fatalf(`Unable to fetch HTTPS TLS key pairs`)
	}

	if len(privateKey.Certificates) == 0 {
		name := proto.GetService().GetName()
		cert, err := tls_.CreateSelfSignedCertificate(privateKey.Key, []string{name}, name)
		if err != nil {
			logger.WithError(errors.WithStack(err)).Fatalf(`Could not generate a self signed TLS certificate`)
		}

		privateKey.Certificates = []*x509.Certificate{cert}
		if err := keyManager.DeleteKey(context.TODO(), tlsKeyName(), privateKey.KeyID); err != nil {
			logger.WithError(errors.WithStack(err)).Fatal(`Could not update (delete) the self signed TLS certificate`)
		}

		if err := keyManager.AddKey(context.TODO(), tlsKeyName(), privateKey); err != nil {
			logger.WithError(errors.WithStack(err)).Fatal(`Could not update (add) the self signed TLS certificate`)
		}
	}

	block, err := key.PEMBlockForKey(privateKey.Key)
	if err != nil {
		logger.WithError(errors.WithStack(err)).Fatalf("Could not encode key to PEM")
	}

	if len(privateKey.Certificates) == 0 {
		logger.Fatal("TLS certificate chain can not be empty")
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: privateKey.Certificates[0].Raw})
	pemKey := pem.EncodeToMemory(block)
	ct, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		logger.WithError(errors.WithStack(err)).Fatal("Could not decode certificate")
	}

	certPool, err = tls_.LoadX509CertificatePool(nil, "", "", ct)
	if err != nil {
		logger.WithError(errors.WithStack(err)).Fatal("TLS certificate pool can not be loaded")
	}

	return []tls.Certificate{ct}, certPool
}

func AsymmetricKeypair(ctx context.Context, g key.KeyGenerator, setId string) (public, private *jose.JSONWebKey, err error) {
	priv, err := GetOrCreateKey(ctx, g, setId, "private")
	if err != nil {
		return nil, nil, err
	}

	pub, err := GetOrCreateKey(ctx, g, setId, "public")
	if err != nil {
		return nil, nil, err
	}

	return pub, priv, nil
}

func GetOrCreateKey(ctx context.Context, g key.KeyGenerator, setId, prefix string) (*jose.JSONWebKey, error) {
	logger := GlobalProvider().Logger().WithField("module", "provider.tls")
	keyManager := GlobalProvider().KeyManager()

	keys_, err := keyManager.GetKeySet(ctx, setId)
	if errors.Cause(err) == orm.ErrNotFound || keys_ != nil && len(keys_.Keys) == 0 {
		logger.Warnf("JSON Web Key Set \"%s\" does not exist yet, generating new key pair...", setId)
		keys_, err = key.CreateKey(ctx, keyManager, g, setId)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	key_, err := key.FindKeyByPrefix(keys_, prefix)
	if err != nil {
		logger.Warnf("JSON Web Key with prefix %s not found in JSON Web Key Set %s, generating new key pair...", prefix, setId)

		keys_, err = key.CreateKey(ctx, keyManager, g, setId)
		if err != nil {
			return nil, err
		}

		key_, err = key.FindKeyByPrefix(keys_, prefix)
		if err != nil {
			return nil, err
		}
	}

	return key_, nil
}
